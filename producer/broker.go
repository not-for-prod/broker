package producer

import (
	"context"

	"github.com/not-for-prod/broker"
)

type Broker interface {
	Push(ctx context.Context, r []broker.Event) error
}
