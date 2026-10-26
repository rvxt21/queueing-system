package domain

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ID          uuid.UUID
	ServiceCode string
	Name        string
	Prefix      string // short code at the start of ticket numbers, e.g. "PN" in "PN-042"
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
