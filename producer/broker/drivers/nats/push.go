package nats

import (
	context "context"
	"fmt"

	"github.com/nats-io/nats.go"
	model "github.com/not-for-prod/broker"
)

func (i *Implementation) Push(_ context.Context, r []model.Event) error {
	for _, e := range r {
		msg := convertEventToMsg(e)

		_, err := i.js.PublishMsg(msg)
		if err != nil {
			return fmt.Errorf("publish to NATS failed: %w", err)
		}
	}

	return nil
}

func convertEventToMsg(e model.Event) *nats.Msg {
	msg := &nats.Msg{
		Subject: e.Topic,
		Header:  nats.Header{},
		Data:    e.Body,
	}

	for k, v := range e.Headers {
		msg.Header.Set(k, v)
	}

	msg.Header.Set("id", e.ID.String())

	return msg
}
