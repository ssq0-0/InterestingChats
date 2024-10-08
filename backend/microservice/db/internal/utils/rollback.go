package utils

import (
	"database/sql"
	"log"
)

// RollbackOnError to perform a transaction rollback if an error occurs.
func RollbackOnError(tx *sql.Tx, err error) {
	if p := recover(); p != nil || err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("failed to rollback transaction: %v", rbErr)
		}
	}
}
