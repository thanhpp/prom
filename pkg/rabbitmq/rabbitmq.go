package rabbitmq

import (
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)

const (
	jsonType = "application/json"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func (r RabbitMQ) pkgError(op string, inErr error) (err error) {
	return fmt.Errorf("Pkg: RabbitMQ. Op: %s. Error: %v", op, err)
}

func (r *RabbitMQ) Setup(serverURL string) (err error) {
	conn, err := amqp.Dial(serverURL)
	if err != nil {
		return r.pkgError("Dial server", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return r.pkgError("Create channel", err)
	}

	r.conn = conn
	r.channel = ch

	return nil
}

func (r *RabbitMQ) CloseConnection() (err error) {
	if r.channel == nil {
		return r.pkgError("Channel check", errors.New("Empty channel"))
	}
	if err = r.channel.Close(); err != nil {
		return r.pkgError("Close channel", err)
	}

	if r.conn == nil {
		return r.pkgError("Connection check", errors.New("Empty connection"))
	}
	if err = r.conn.Close(); err != nil {
		return r.pkgError("Close connection", err)
	}

	return nil
}

func (r *RabbitMQ) CreateQueue(queueName string) (err error) {
	if r.channel == nil {
		return r.pkgError("Channel check", errors.New("Empty channel"))
	}

	_, err = r.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return r.pkgError("Create queue", err)
	}

	return nil
}

func (r *RabbitMQ) PublishMsg(queue string, message []byte) (err error) {
	if r.channel == nil {
		return r.pkgError("Channel check", errors.New("Empty channel"))
	}

	msg := amqp.Publishing{
		ContentType: jsonType,
		Body:        message,
	}

	if err = r.channel.Publish("", queue, false, false, msg); err != nil {
		return r.pkgError("Publish message", err)
	}

	return nil
}

func (r *RabbitMQ) GetConsumerChan(queue string) (messages <-chan amqp.Delivery, err error) {
	if r.channel == nil {
		return nil, r.pkgError("Channel check", errors.New("Empty channel"))
	}

	messages, err = r.channel.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		return nil, r.pkgError("Create consume channel", err)
	}

	return messages, nil
}
