package domain

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ID          uuid.UUID
	ServiceCode string
	Name        string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
