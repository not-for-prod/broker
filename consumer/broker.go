package consumer

import (
	"context"

	"github.com/not-for-prod/broker/models"
)

type Broker interface {
	Pull(ctx context.Context, batchSize uint64) ([]models.Event, error)
	Commit(ctx context.Context, events []models.Event) error
}
