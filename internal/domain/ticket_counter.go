package domain

import "time"

type TicketCounter struct {
	ServiceCode string
	Day         time.Time
	LastNumber  int32
}

func (c *TicketCounter) Next() int32 {
	c.LastNumber++

	return c.LastNumber
}
