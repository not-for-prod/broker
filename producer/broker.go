package producer

import (
	"context"

	"github.com/not-for-prod/broker"
)

//go:generate moq -out brocker_mock.go . Broker

type Broker interface {
	Push(ctx context.Context, r []broker.Event) error
}
