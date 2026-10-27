package broker

import "errors"

var (
	ErrPublishNotConfirmed = errors.New("rabbitmq did not confirm the message")
	ErrOpenChannel         = errors.New("could not open a channel for a consumer")
)
