package node

import (
	"time"
)

const (
	onlineScoreRatio = 100.0
)

func (m *Manager) getScale() map[string]struct {
	min, max, codeNum int
} {
	return map[string]struct {
		min, max, codeNum int
	}{
		"A": {90, 100, 3},
		"B": {50, 89, 2},
		"C": {0, 49, 1},
	}
}

func (m *Manager) getSelectCodeNum(nodeID string) int {
	score := int(m.getNodeScore(nodeID))

	for _, rangeScore := range m.getScale() {
		if score >= rangeScore.min && score <= rangeScore.max {
			return rangeScore.codeNum
		}
	}

	return 0
}

func (m *Manager) getNodeScore(nodeID string) float64 {
	// online time
	info, err := m.LoadNodeInfo(nodeID)
	if err != nil {
		log.Errorf("LoadNodeInfo err:%s", err.Error())
		return 0
	}

	minutes := time.Now().Sub(info.FirstTime).Minutes()
	onlineRatio := float64(info.OnlineDuration) / minutes
	if onlineRatio > 1 {
		onlineRatio = 1
	}

	return onlineScoreRatio * onlineRatio
}
