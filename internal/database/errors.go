package database

import "errors"

var (
	ErrOpenPostgresConnection = errors.New("error open postgres connection")
	ErrPingingDatabase        = errors.New("error pinging database")
	ErrGetServices            = errors.New("error getting services")
	ErrGetServiceByCode       = errors.New("error getting service by code")
	ErrCreateTicket           = errors.New("error creating ticket")
	ErrLockTicketCounter      = errors.New("error locking ticket counter")
	ErrSaveTicketCounter      = errors.New("error saving ticket counter")
	ErrNoTicketToCreate       = errors.New("createFn returned no ticket")
	ErrBeginTx                = errors.New("error beginning transaction")
	ErrCommitTx               = errors.New("error committing transaction")
)
