package broker

import (
	"context"
	"errors"
	"fmt"

	"github.com/rvxt21/queueing-system/internal/domain"
)

func (r *RabbitMQ) Consume(ctx context.Context, queue, consumerTag string) (<-chan domain.Delivery, error) {
	ch, err := r.currentConn().Channel()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrOpenChannel, err)
	}

	err = ch.Qos(1, 0, false)
	if err != nil {
		return nil, fmt.Errorf("set prefetch: %w", errors.Join(err, ch.Close()))
	}

	msgs, err := ch.Consume(
		queue,
		consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("start consumer on %s: %w", queue, errors.Join(err, ch.Close()))
	}

	out := make(chan domain.Delivery)

	go func() {
		defer close(out)
		defer func() { _ = ch.Close() }()

		for {
			select {
			case <-ctx.Done():
				_ = ch.Cancel(consumerTag, false)
				return
			case d, ok := <-msgs:
				if !ok {
					return
				}

				msg := domain.Delivery{
					MessageID: d.MessageId,
					Body:      d.Body,
					Ack:       func() error { return d.Ack(false) },
					Nack:      func(requeue bool) error { return d.Nack(false, requeue) },
				}

				select {
				case <-ctx.Done():
					_ = d.Nack(false, true)
					return
				case out <- msg:
				}
			}
		}
	}()

	return out, nil
}
