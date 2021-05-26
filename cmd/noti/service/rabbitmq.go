package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/thanhpp/prom/cmd/noti/repository"
	"github.com/thanhpp/prom/cmd/noti/repository/entity"
	"github.com/thanhpp/prom/pkg/booting"
	"github.com/thanhpp/prom/pkg/logger"
	"github.com/thanhpp/prom/pkg/rabbitmq"
)

type RabbitMQService struct {
	srv *rabbitmq.RabbitMQ
}

func (r *RabbitMQService) srvError(op string, err error) error {
	return fmt.Errorf("Srv: RabbitMQ. Op: %s. Err: %v", op, err)
}

func (r *RabbitMQService) Connect(serverURL string) (err error) {
	r.srv = new(rabbitmq.RabbitMQ)

	if err = r.srv.Setup(serverURL); err != nil {
		return r.srvError("Setup", err)
	}

	return nil
}

func (r *RabbitMQService) CreateMsgDaemon(ctx context.Context) (daemon booting.Daemon, err error) {
	if r.srv == nil {
		return nil, errors.New("Empty service")
	}

	if err = r.srv.CreateQueue(rabbitmq.NotificationQueue); err != nil {
		return nil, err
	}

	daemon = func(ctx context.Context) (start func() error, stop func()) {
		start = func() error {
			msgq, err := r.srv.GetConsumerChan(rabbitmq.NotificationQueue)
			if err != nil {
				return r.srvError("Get consumer queue", err)
			}

			for {
				select {
				case msg := <-msgq:
					if err := r.handleCreateMsg(ctx, msg.Body); err != nil {
						logger.Get().Errorf("Handle message error: %v", err)
					}
				case <-ctx.Done():
					return r.srvError("Context done", ctx.Err())
				}
			}
		}

		stop = func() {
			if err := r.srv.CloseConnection(); err != nil {
				logger.Get().Errorf("Stop rabbitMQ error: %v", err)
			}
		}
		return start, stop
	}
	return daemon, nil
}

func (r *RabbitMQService) handleCreateMsg(ctx context.Context, msg []byte) (err error) {
	if ctx.Err() != nil {
		return r.srvError("context check", err)
	}

	if msg == nil {
		return nil
	}

	newNotiMsg := new(rabbitmq.NewNotificationMsg)
	if err = json.Unmarshal(msg, newNotiMsg); err != nil {
		return r.srvError("Unmarshal new noti msg", err)
	}

	noti := &entity.Notification{
		CardID:  newNotiMsg.CardID,
		Seen:    false,
		Content: newNotiMsg.Content,
	}

	if err = repository.Get().CreateNotification(ctx, noti, newNotiMsg.UserIDs); err != nil {
		return err
	}

	return nil
}
