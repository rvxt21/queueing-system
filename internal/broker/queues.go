package broker

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	QueueIdentity  = "tickets.identity"
	QueueTax       = "tickets.tax"
	QueueTransport = "tickets.transport"
	QueueCivil     = "tickets.civil"
	QueueResidence = "tickets.residence"
)

var ticketQueues = []string{
	QueueIdentity,
	QueueTax,
	QueueTransport,
	QueueCivil,
	QueueResidence,
}

var serviceQueues = map[string]string{
	// identity documents
	"passport_new":     QueueIdentity,
	"passport_collect": QueueIdentity,
	"id_card_new":      QueueIdentity,
	"id_card_collect":  QueueIdentity,

	// tax and business
	"tax_registration":      QueueTax,
	"tax_declaration":       QueueTax,
	"tax_certificate":       QueueTax,
	"business_registration": QueueTax,

	// vehicles and driving
	"vehicle_registration":    QueueTransport,
	"driving_license_new":     QueueTransport,
	"driving_license_replace": QueueTransport,
	"parking_permit":          QueueTransport,

	// civil registry
	"birth_certificate":     QueueCivil,
	"marriage_certificate":  QueueCivil,
	"marriage_registration": QueueCivil,
	"criminal_record":       QueueCivil,

	// residence, property and general inquiries
	"address_registration":   QueueResidence,
	"address_deregistration": QueueResidence,
	"building_permit":        QueueResidence,
	"complaint":              QueueResidence,
}

func declareQueues(ch *amqp.Channel) error {
	for _, name := range ticketQueues {
		_, err := ch.QueueDeclare(
			name,  // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			return fmt.Errorf("declare queue %s: %w", name, err)
		}
	}

	return nil
}

func bindQueues(ch *amqp.Channel, exchange string) error {
	for serviceCode, queueName := range serviceQueues {
		err := ch.QueueBind(
			queueName,
			serviceCode,
			exchange,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("bind %s to queue %s: %w", serviceCode, queueName, err)
		}
	}

	return nil
}
