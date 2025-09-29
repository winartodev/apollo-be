package database

import (
	"database/sql"
	"fmt"
)

const (
	errorFailedToCloseStatement    = "failed to close statement: %v"
	errorAlsoFiledToCloseStatement = "%v, also failed to close statement: %v"
	ErrFailedPrepareStatement      = "failed to prepare statement: %v"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase(db *sql.DB) (*Database, error) {
	return &Database{
		DB: db,
	}, nil
}

func (d *Database) CloseStatement(stmt *sql.Stmt, errPtr *error) {
	if closeErr := stmt.Close(); closeErr != nil {
		if *errPtr == nil {
			*errPtr = fmt.Errorf(errorFailedToCloseStatement, closeErr)
		} else {
			*errPtr = fmt.Errorf(errorAlsoFiledToCloseStatement, *errPtr, closeErr)
		}
	}
}
