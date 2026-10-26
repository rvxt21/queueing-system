package broker

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rvxt21/queueing-system/internal/domain"
)

func (r *RabbitMQ) Publish(ctx context.Context, msg *domain.Outbox) error {
	ctx, cancel := context.WithTimeout(ctx, r.config.PublishTimeout)
	defer cancel()

	confirmation, err := r.ch.PublishWithDeferredConfirmWithContext(ctx,
		r.config.Exchange, // exchange
		msg.RoutingKey,    // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			MessageId:    msg.ID.String(),
			Type:         msg.EventType,
			Timestamp:    msg.CreatedAt,
			Body:         []byte(msg.Payload),
		})
	if err != nil {
		return err
	}

	acked, err := confirmation.WaitContext(ctx)
	if err != nil {
		return err
	}

	if !acked {
		return fmt.Errorf("%w: message %s", ErrPublishNotConfirmed, msg.ID)
	}

	return nil
}
