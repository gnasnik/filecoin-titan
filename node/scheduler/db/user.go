package db

import (
	"fmt"
)

// SaveAssetUser save asset and user info
func (n *SQLDB) SaveAssetUser(hash, userID string) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (hash, user_id, created_time) 
		        VALUES (?, ?, NOW()) 
				ON DUPLICATE KEY UPDATE hash=?, user_id=?, created_time=NOW() `, userAssetTable)
	_, err := n.db.Exec(query, hash, userID, hash, userID)

	return err
}

// DeleteAssetUser delete asset and user info
func (n *SQLDB) DeleteAssetUser(hash, userID string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE hash=? AND user_id=? `, userAssetTable)
	_, err := n.db.Exec(query, hash, userID)

	return err
}
