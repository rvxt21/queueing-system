package database

import "errors"

var (
	ErrOpenPostgresConnection = errors.New("error open postgres connection")
	ErrPingingDatabase        = errors.New("error pinging database")
)
