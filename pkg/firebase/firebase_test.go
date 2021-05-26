package firebase_test

import (
	"context"
	"testing"

	"github.com/thanhpp/prom/pkg/firebase"
)

func createFirebaseClient(ctx context.Context, keyPath string) (f *firebase.FireBase, err error) {
	f = new(firebase.FireBase)
	if err := f.CreateClient(ctx, keyPath); err != nil {
		return nil, err
	}

	return f, nil
}

func TestCreateClient(t *testing.T) {
	var (
		keyPath = "/home/thanhpp/go/src/github.com/thanhpp/prom/prom-eef45-firebase-adminsdk-5f4bn-b133193339.json"
		ctx     = context.Background()
	)

	if _, err := createFirebaseClient(ctx, keyPath); err != nil {
		t.Error(err)
		return
	}
}

func TestSendMsgToTopic(t *testing.T) {
	var (
		keyPath = "/home/thanhpp/go/src/github.com/thanhpp/prom/prom-eef45-firebase-adminsdk-5f4bn-b133193339.json"
		ctx     = context.Background()
		topic   = "test"
		data    = map[string]string{
			"hello": "world",
		}
	)

	f, err := createFirebaseClient(ctx, keyPath)
	if err != nil {
		t.Error(err)
		return
	}

	if err := f.SendMsgToTopic(ctx, topic, data); err != nil {
		t.Error(err)
		return
	}
}
