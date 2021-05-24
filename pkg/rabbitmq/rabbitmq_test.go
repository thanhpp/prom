package rabbitmq_test

import (
	"testing"

	"github.com/thanhpp/prom/pkg/rabbitmq"
)

func TestSetup(t *testing.T) {
	var (
		serverURL = "amqp://guest:guest@localhost:5672/"
		rmq       = new(rabbitmq.RabbitMQ)
	)

	if err := rmq.Setup(serverURL); err != nil {
		t.Error(err)
		return
	}

	if err := rmq.CloseConnection(); err != nil {
		t.Error(err)
		return
	}
}

func TestMessage(t *testing.T) {
	var (
		serverURL = "amqp://guest:guest@localhost:5672/"
		rmq       = new(rabbitmq.RabbitMQ)
		queueName = "testqueue"
	)

	if err := rmq.Setup(serverURL); err != nil {
		t.Error(err)
		return
	}

	if err := rmq.CreateQueue(queueName); err != nil {
		t.Error(err)
		return
	}

	if err := rmq.PublishMsg(queueName, []byte("message")); err != nil {
		t.Error(err)
		return
	}

	msgs, err := rmq.GetConsumerChan(queueName)
	if err != nil {
		t.Error(err)
		return
	}

	for msg := range msgs {
		t.Log(msg.Body)
	}

	if err := rmq.CloseConnection(); err != nil {
		t.Error(err)
		return
	}
}
