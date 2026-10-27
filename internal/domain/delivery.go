package domain

type Delivery struct {
	MessageID string
	Body      []byte
	Ack       func() error
	Nack      func(requeue bool) error
}
