package nats

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	model "github.com/not-for-prod/broker/models"
)

func (i *Implementation) Pull(ctx context.Context, batchSize uint64) ([]model.Event, error) {
	msgs, err := i.sub.Fetch(int(batchSize), nats.Context(ctx))
	if err != nil && err != context.Canceled {
		return nil, fmt.Errorf("fetch messages failed: %w", err)
	}

	events := make([]model.Event, 0, len(msgs))

	i.bufferMu.Lock()
	defer i.bufferMu.Unlock()

	for _, msg := range msgs {
		id := msg.Header.Get("id")

		event := model.Event{
			Ctx:       ctx,
			ID:        model.EventID(id),
			Topic:     msg.Subject,
			Partition: "",
			Headers:   make(map[string]string),
			Body:      msg.Data,
		}

		for k, v := range msg.Header {
			event.Headers[k] = v[0]
		}

		events = append(events, event)
		i.buffer[event.ID] = msg
	}

	return events, nil
}
