package consumer

import (
	"context"

	"github.com/not-for-prod/broker"
)

//go:generate moq -out brocker_mock.go . Broker

type Broker interface {
	Pull(ctx context.Context, batchSize uint64) ([]broker.Event, error)
	Commit(ctx context.Context, events []broker.Event) error
}
