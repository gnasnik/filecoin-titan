package workload

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/Filecoin-Titan/titan/api/types"
	"github.com/Filecoin-Titan/titan/node/modules/dtypes"
	"github.com/Filecoin-Titan/titan/node/scheduler/db"
	"github.com/Filecoin-Titan/titan/node/scheduler/leadership"
	"github.com/Filecoin-Titan/titan/node/scheduler/node"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"
)

var log = logging.Logger("workload")

const (
	oneDay = 24 * time.Hour

	// Process 1000 pieces of workload result data at a time
	vWorkloadLimit = 1000
)

// Manager node workload
type Manager struct {
	config        dtypes.GetSchedulerConfigFunc
	leadershipMgr *leadership.Manager
	nodeMgr       *node.Manager
	*db.SQLDB

	profit float64
}

// NewManager return new node manager instance
func NewManager(sdb *db.SQLDB, configFunc dtypes.GetSchedulerConfigFunc, lmgr *leadership.Manager, nmgr *node.Manager) *Manager {
	manager := &Manager{
		config:        configFunc,
		leadershipMgr: lmgr,
		SQLDB:         sdb,
		nodeMgr:       nmgr,
	}

	go manager.startHandleWorkloadResult()

	return manager
}

func (m *Manager) startHandleWorkloadResult() {
	now := time.Now()

	nextTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if now.After(nextTime) {
		nextTime = nextTime.Add(oneDay)
	}

	duration := nextTime.Sub(now)

	timer := time.NewTimer(duration)
	defer timer.Stop()

	for {
		<-timer.C

		log.Debugln("start workload timer...")
		m.handleWorkloadResult()

		timer.Reset(oneDay)
	}
}

func (m *Manager) handleWorkloadResult() {
	if !m.leadershipMgr.RequestAndBecomeMaster() {
		return
	}

	defer func() {
		log.Infoln("handleWorkloadResult end")
	}()

	log.Infoln("handleWorkloadResult start")

	profit := m.getValidationProfit()

	// do handle workload result
	for {
		rows, err := m.LoadUnprocessedWorkloadResults(vWorkloadLimit)
		if err != nil {
			log.Errorf("LoadWorkloadResults err:%s", err.Error())
			return
		}

		ids := make(map[string]types.WorkloadStatus)
		nodeProfits := make(map[string]float64)

		for rows.Next() {
			record := &types.WorkloadRecord{}
			err = rows.StructScan(record)
			if err != nil {
				log.Errorf("ValidationResultInfo StructScan err: %s", err.Error())
				continue
			}

			// check workload ...
			status := m.checkWorkload(record)

			ids[record.ID] = status
			if status == types.WorkloadStatusSucceeded {
				nodeProfits[record.NodeID] += profit
			}
		}
		rows.Close()

		if len(ids) == 0 {
			return
		}

		err = m.UpdateNodeProfitsByWorkloadResult(ids, nodeProfits)
		if err != nil {
			log.Errorf("UpdateNodeProfitsByWorkloadResult err:%s", err.Error())
		}
	}
}

// get the profit of validation
func (m *Manager) getValidationProfit() float64 {
	cfg, err := m.config()
	if err != nil {
		log.Errorf("get config err:%s", err.Error())
		return 0
	}

	return cfg.WorkloadProfit
}

func (m *Manager) checkWorkload(record *types.WorkloadRecord) types.WorkloadStatus {
	nWorkload := &types.Workload{}
	if len(record.NodeWorkload) > 0 {
		dec := gob.NewDecoder(bytes.NewBuffer(record.NodeWorkload))
		err := dec.Decode(nWorkload)
		if err != nil {
			log.Errorf("decode data to *types.Workload error: %w", err)
			return types.WorkloadStatusFailed
		}
	}

	cWorkload := &types.Workload{}
	if len(record.ClientWorkload) > 0 {
		dec := gob.NewDecoder(bytes.NewBuffer(record.ClientWorkload))
		err := dec.Decode(cWorkload)
		if err != nil {
			log.Errorf("decode data to *types.Workload error: %w", err)
			return types.WorkloadStatusFailed
		}
	}

	if nWorkload.DownloadSize != cWorkload.DownloadSize {
		return types.WorkloadStatusFailed
	}

	// TODO other ...

	return types.WorkloadStatusSucceeded
}

// HandleUserWorkload handle user workload
func (m *Manager) HandleUserWorkload(data []byte, node *node.Node) error {
	reports := make([]*types.WorkloadReport, 0)
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	err := dec.Decode(&reports)
	if err != nil {
		return xerrors.Errorf("decode data to []*types.WorkloadReport error: %w", err)
	}

	size := int64(0)
	// TODO merge workload report by token
	for _, rp := range reports {
		if rp.Workload == nil {
			log.Errorf("workload cannot empty %#v", *rp)
			continue
		}

		size += rp.Workload.DownloadSize

		// replace clientID with nodeID
		if node != nil {
			rp.ClientID = node.NodeID
		}
		if err = m.handleWorkloadReport(rp.NodeID, rp, true); err != nil {
			log.Errorf("handler user workload report error %s, token id %s", err.Error(), rp.TokenID)
			continue
		}

		size += rp.Workload.DownloadSize
	}

	if node != nil {
		node.DownloadTraffic += size
	}

	return nil
}

// HandleNodeWorkload handle node workload
func (m *Manager) HandleNodeWorkload(data []byte, node *node.Node) error {
	reports := make([]*types.WorkloadReport, 0)
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	err := dec.Decode(&reports)
	if err != nil {
		return xerrors.Errorf("decode data to []*types.WorkloadReport error: %w", err)
	}

	// TODO merge workload report by token
	for _, rp := range reports {
		if err = m.handleWorkloadReport(node.NodeID, rp, false); err != nil {
			log.Errorf("handler node workload report error %s", err.Error())
			continue
		}

		node.UploadTraffic += rp.Workload.DownloadSize
	}

	return nil
}

func (m *Manager) handleWorkloadReport(nodeID string, report *types.WorkloadReport, isClient bool) error {
	workloadRecord, err := m.LoadWorkloadRecord(report.TokenID)
	if err != nil {
		return xerrors.Errorf("load token payload and workloads with token id %s error: %w", report.TokenID, err)
	}
	if isClient && workloadRecord.NodeID != nodeID {
		return fmt.Errorf("token payload node id %s, but report node id is %s", workloadRecord.NodeID, report.NodeID)
	}

	if !isClient && workloadRecord.ClientID != report.ClientID {
		return fmt.Errorf("token payload client id %s, but report client id is %s", workloadRecord.ClientID, report.ClientID)
	}

	if workloadRecord.Expiration.Before(time.Now()) {
		return fmt.Errorf("token payload expiration %s < %s", workloadRecord.Expiration.Local().String(), time.Now().Local().String())
	}

	workloadBytes := workloadRecord.NodeWorkload
	if isClient {
		workloadBytes = workloadRecord.ClientWorkload
	}

	workload := &types.Workload{}
	if len(workloadBytes) > 0 {
		dec := gob.NewDecoder(bytes.NewBuffer(workloadBytes))
		err = dec.Decode(workload)
		if err != nil {
			return xerrors.Errorf("decode data to []*types.Workload error: %w", err)
		}
	}

	workload = m.mergeWorkloads([]*types.Workload{workload, report.Workload})

	buffer := &bytes.Buffer{}
	enc := gob.NewEncoder(buffer)
	err = enc.Encode(workload)
	if err != nil {
		return xerrors.Errorf("encode data to Buffer error: %w", err)
	}

	if isClient {
		workloadRecord.ClientWorkload = buffer.Bytes()
	} else {
		workloadRecord.NodeWorkload = buffer.Bytes()
	}

	return m.UpdateWorkloadRecord(workloadRecord)
}

func (m *Manager) mergeWorkloads(workloads []*types.Workload) *types.Workload {
	if len(workloads) == 0 {
		return nil
	}

	// costTime := int64(0)
	downloadSize := int64(0)
	startTime := int64(0)
	endTime := int64(0)
	speedCount := int64(0)
	accumulateSpeed := int64(0)

	for _, workload := range workloads {
		if workload.DownloadSpeed > 0 {
			accumulateSpeed += workload.DownloadSpeed
			speedCount++
		}
		downloadSize += workload.DownloadSize

		if startTime == 0 || workload.StartTime < startTime {
			startTime = workload.StartTime
		}

		if workload.EndTime > endTime {
			endTime = workload.EndTime
		}
	}

	downloadSpeed := int64(0)
	if speedCount > 0 {
		downloadSpeed = accumulateSpeed / speedCount
	}
	return &types.Workload{DownloadSpeed: downloadSpeed, DownloadSize: downloadSize, StartTime: startTime, EndTime: endTime}
}
