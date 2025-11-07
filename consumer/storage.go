package consumer

import (
	"context"

	"github.com/not-for-prod/broker"
)

//go:generate moq -out storage_mock.go . Storage TxManager

type Storage interface {
	SetNX(ctx context.Context, e []broker.Event) ([]bool, error)
}

type TxManager interface {
	Do(ctx context.Context, fn func(context.Context) error) error
}
