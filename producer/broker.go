package producer

import (
	"context"

	"github.com/not-for-prod/broker/models"
)

type Broker interface {
	Push(ctx context.Context, r []models.Event) error
}
