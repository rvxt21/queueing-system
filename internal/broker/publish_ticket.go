package broker

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rvxt21/queueing-system/internal/domain"
)

func (r *RabbitMQ) PublishTicket(ctx context.Context, ticket *domain.Ticket) error {
	ctx, cancel := context.WithTimeout(ctx, r.config.PublishTimeout)
	defer cancel()

	data, err := json.Marshal(newTicketMessage(ticket))
	if err != nil {
		return fmt.Errorf("marshal ticket %s: %w", ticket.ID, err)
	}

	confirmation, err := r.ch.PublishWithDeferredConfirmWithContext(ctx,
		r.config.Exchange,  // exchange
		ticket.ServiceCode, // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			MessageId:    ticket.ID.String(),
			Body:         data,
		})
	if err != nil {
		return err
	}

	acked, err := confirmation.WaitContext(ctx)
	if err != nil {
		return err
	}

	if !acked {
		return fmt.Errorf("%w: ticket %s", ErrPublishNotConfirmed, ticket.ID)
	}

	return nil
}
