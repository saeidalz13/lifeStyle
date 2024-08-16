package dberr

import "errors"

var (
	ErrNoMigration       = errors.New("no migration")
	ErrDirtyDb           = errors.New("database is dirty")
	ErrFileNotExists     = errors.New("first .: file does not exist")
	ErrMigrationNoChange = errors.New("no change")
)
