package domain

import "errors"

var (
	ErrServiceNotFound = errors.New("service not found")
	ErrTicketNotFound  = errors.New("ticket not found")
	ErrServiceInactive = errors.New("service is not active")
)
