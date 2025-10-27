package consumer

import (
	"context"

	"github.com/not-for-prod/broker/models"
)

type Storage interface {
	SetNX(ctx context.Context, e []models.Event) ([]bool, error)
}

type TxManager interface {
	Do(ctx context.Context, fn func(context.Context) error) error
}
