package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	connection *amqp.Connection
}

func (p *Producer) setup() error {
	channel, err := p.connection.Channel()

	if err != nil {
		return err
	}

	defer channel.Close()

	return declareExchange(channel)
}

// Push adds events to RabbitMQ
func (p *Producer) Push(event string, severity string) error {
	// open channel for bulking of messages
	channel, err := p.connection.Channel()

	if err != nil {
		return err
	}

	// close channel
	defer channel.Close()

	// sends message to the RabbitMq
	err = channel.Publish(
		"notes_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)

	if err != nil {
		return err
	}
	log.Println(err)
	return nil
}

// NewProducer sets up a producer.
func NewProducer(conn *amqp.Connection) (Producer, error) {
	producer := Producer{
		connection: conn,
	}

	err := producer.setup()
	if err != nil {
		return Producer{}, err
	}

	return producer, nil
}
