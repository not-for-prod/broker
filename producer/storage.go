package producer

import (
	"context"

	"github.com/not-for-prod/broker/models"
)

type Storage interface {
	Push(ctx context.Context, e []models.Event) error
	GetOffset(ctx context.Context, producerName string) (uint64, error)
	CommitOffset(ctx context.Context, producerName string, offset uint64) error
	ListRecords(ctx context.Context, limit, offset uint64) ([]models.Event, error)
}

type TxManager interface {
	Do(ctx context.Context, fn func(context.Context) error) error
}
