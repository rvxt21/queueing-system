package broker

import "errors"

var ErrPublishNotConfirmed = errors.New("rabbitmq did not confirm the message")
