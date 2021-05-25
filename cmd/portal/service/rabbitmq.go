package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/thanhpp/prom/pkg/rabbitmq"
)

type RabbitMQService struct {
	srv *rabbitmq.RabbitMQ
}

var rmq = new(RabbitMQService)

func (r *RabbitMQService) srvError(op string, err error) error {
	return fmt.Errorf("Srv: RabbitMQ. Op: %s. Err: %v", op, err)
}

func SetRabbitMQ(serverURL string) (err error) {
	rmq.srv = new(rabbitmq.RabbitMQ)
	if err = rmq.srv.Setup(serverURL); err != nil {
		return rmq.srvError("Setup", err)
	}

	if err = rmq.srv.CreateQueue(rabbitmq.NotificationQueue); err != nil {
		return rmq.srvError("Create new queue", err)
	}

	return nil
}

func GetRabbitMQ() *RabbitMQService {
	return rmq
}

func (r *RabbitMQService) SendNewNotiMsg(msg *rabbitmq.NewNotificationMsg) (err error) {
	if r.srv == nil {
		return errors.New("Empty service")
	}

	msgByte, err := json.Marshal(msg)
	if err != nil {
		return r.srvError("Marshal msg", err)
	}

	if err = r.srv.PublishMsg(rabbitmq.NotificationQueue, msgByte); err != nil {
		return r.srvError("Publish msg", err)
	}

	return nil
}

func (r *RabbitMQService) Close() (err error) {
	if r.srv == nil {
		return errors.New("Empty service")
	}

	if err = r.srv.CloseConnection(); err != nil {
		return err
	}

	return nil
}
