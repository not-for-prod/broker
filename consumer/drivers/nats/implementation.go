package nats

import (
	"fmt"
	"sync"

	"github.com/nats-io/nats.go"
	model "github.com/not-for-prod/broker/models"
)

type Implementation struct {
	sub *nats.Subscription

	bufferMu sync.Mutex
	buffer   map[model.EventID]*nats.Msg
}

func NewImplementation(js nats.JetStream, topic, consumerGroup string) (*Implementation, error) {
	sub, err := js.PullSubscribe(topic, consumerGroup)
	if err != nil {
		return nil, fmt.Errorf("create pull subscription failed: %w", err)
	}

	return &Implementation{
		sub:    sub,
		buffer: make(map[model.EventID]*nats.Msg),
	}, nil
}
