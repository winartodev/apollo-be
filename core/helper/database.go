package helper

import (
	"database/sql"
	"fmt"
)

const (
	errorFailedToCloseStatement    = "failed to close statement: %v"
	errorAlsoFiledToCloseStatement = "%v, also failed to close statement: %v"
)

type DatabaseUtil struct {
	DB *sql.DB
}

func NewDatabaseUtil(db *sql.DB) (*DatabaseUtil, error) {
	return &DatabaseUtil{
		DB: db,
	}, nil
}

func (d *DatabaseUtil) CloseStatement(stmt *sql.Stmt, errPtr *error) {
	if closeErr := stmt.Close(); closeErr != nil {
		if *errPtr == nil {
			*errPtr = fmt.Errorf(errorFailedToCloseStatement, closeErr)
		} else {
			*errPtr = fmt.Errorf(errorAlsoFiledToCloseStatement, *errPtr, closeErr)
		}
	}
}
