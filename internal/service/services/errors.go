package services

import "errors"

var (
	ErrListServices   = errors.New("error listing services")
	ErrRequestService = errors.New("error requesting service")
)
