package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type FireBase struct {
	client *messaging.Client
}

func (f *FireBase) pkgError(op string, err error) error {
	return fmt.Errorf("Pkg: Firebase. Op: %s. Err: %v", op, err)
}

func (f *FireBase) CreateClient(ctx context.Context, keyPath string) (err error) {
	opt := option.WithCredentialsFile(keyPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return f.pkgError("Create app", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return f.pkgError("Create client", err)
	}

	f.client = client

	return nil
}

func (f *FireBase) SendMsgToTopic(ctx context.Context, topic string, data map[string]string) (err error) {
	message := &messaging.Message{
		Data:  data,
		Topic: topic,
	}

	// Send a message to the devices subscribed to the provided topic.
	response, err := f.client.Send(ctx, message)
	if err != nil {
		return f.pkgError("Send msg", err)
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)

	return nil
}

func (f *FireBase) SendMsgToIDs(ctx context.Context, tokens []string, data map[string]string) (err error) {
	message := &messaging.MulticastMessage{
		Tokens: tokens,
		Data:   data,
	}

	br, err := f.client.SendMulticast(ctx, message)
	if err != nil {
		return f.pkgError("Send Multicast", err)
	}

	fmt.Println(br)

	return nil
}
