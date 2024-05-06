package db

import (
	"database/sql"
	"fmt"

	"github.com/Filecoin-Titan/titan/api/types"
	"golang.org/x/xerrors"
)

// SaveAWSData save data
func (n *SQLDB) SaveAWSData(infos []types.AWSDataInfo) error {
	tx, err := n.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		err = tx.Rollback()
		if err != nil && err != sql.ErrTxDone {
			log.Errorf("SaveAWSData Rollback err:%s", err.Error())
		}
	}()

	for _, info := range infos {
		if info.Size <= 0 {
			return xerrors.Errorf("%s SaveAWSData size %.2f ", info.Bucket, info.Size)
		}

		sqlString := fmt.Sprintf(`INSERT INTO %s (bucket, replicas, cid, size) VALUES (:bucket, :replicas, :cid, :size) `, awsDataTable)
		tx.NamedExec(sqlString, info)

	}
	return tx.Commit()
}

// UpdateAWSData update aws info
func (n *SQLDB) UpdateAWSData(info *types.AWSDataInfo) error {
	query := fmt.Sprintf(`UPDATE %s SET cid=?, is_distribute=?, distribute_time=NOW() WHERE bucket=?`, awsDataTable)
	_, err := n.db.Exec(query, info.Cid, info.IsDistribute, info.Bucket)
	return err
}

// ListAWSData
func (n *SQLDB) ListAWSData(limit, offset int, isDistribute bool) ([]*types.AWSDataInfo, error) {
	var infos []*types.AWSDataInfo
	query := fmt.Sprintf("SELECT * FROM %s WHERE is_distribute=? LIMIT ? OFFSET ?", awsDataTable)
	err := n.db.Select(&infos, query, isDistribute, limit, offset)
	if err != nil {
		return nil, err
	}

	return infos, nil
}

// LoadAWSData
func (n *SQLDB) LoadAWSData(bucket string) (*types.AWSDataInfo, error) {
	var info types.AWSDataInfo
	query := fmt.Sprintf("SELECT * FROM %s WHERE bucket=? ", awsDataTable)
	err := n.db.Get(&info, query, bucket)
	if err != nil {
		return nil, err
	}

	return &info, nil
}
