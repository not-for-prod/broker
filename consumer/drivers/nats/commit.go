package nats

import (
	context "context"

	model "github.com/not-for-prod/broker/models"
)

func (i *Implementation) CommitOffset(ctx context.Context, events []model.Event) error {
	i.bufferMu.Lock()
	defer i.bufferMu.Unlock()

	for _, event := range events {
		msg, ok := i.buffer[event.ID]
		if !ok {
			continue
		}

		err := msg.AckSync()
		if err != nil {
			return err
		}

		delete(i.buffer, event.ID)
	}

	return nil
}
