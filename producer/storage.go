package producer

import (
	"context"

	"github.com/not-for-prod/broker"
)

//go:generate moq -out storage_mock.go . Storage TxManager

type Storage interface {
	Push(ctx context.Context, e []broker.Event) error
	GetOffset(ctx context.Context, producerName string) (uint64, error)
	CommitOffset(ctx context.Context, producerName string, offset uint64) error
	ListRecords(ctx context.Context, limit, offset uint64) ([]broker.Event, error)
}

type TxManager interface {
	Do(ctx context.Context, fn func(context.Context) error) error
}
