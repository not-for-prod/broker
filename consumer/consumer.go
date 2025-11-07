package consumer

import (
	"context"
	"time"

	"github.com/avast/retry-go"
	"github.com/not-for-prod/broker"
)

type Job func(ctx context.Context, events []broker.Event) error

type Consumer struct {
	Broker    Broker
	Storage   Storage
	TxManager TxManager
	Job       Job
	Logger    broker.Logger
	options   options
	stop      chan struct{}
}

func New(
	broker Broker,
	storage Storage,
	txManager TxManager,
	job Job,
	logger broker.Logger,
	opts ...Option,
) *Consumer {
	options := defaultOptions
	for _, opt := range opts {
		opt.apply(&options)
	}

	return &Consumer{
		Broker:    broker,
		Storage:   storage,
		TxManager: txManager,
		Logger:    logger,
		Job:       job,
		options:   options,
		stop:      make(chan struct{}),
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	ticker := time.NewTicker(c.options.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c.stop:
			return nil
		case <-ticker.C:
			if err := c.consume(ctx); err != nil {
				c.Logger.Error("failed to consume events", "error", err.Error())
			}
		}
	}
}

func (c *Consumer) Stop(_ context.Context) error {
	close(c.stop)
	return nil
}

func (c *Consumer) consume(ctx context.Context) error {
	records, err := c.Broker.Pull(ctx, c.options.batchSize)
	if err != nil {
		return err
	}

	return c.TxManager.Do(
		ctx, func(ctx context.Context) error {
			nx, err := c.Storage.SetNX(ctx, records)
			if err != nil {
				return err
			}

			filtered := make([]broker.Event, 0, len(records))

			for i := range records {
				if nx[i] {
					filtered = append(filtered, records[i])
				}
			}

			err = retry.Do(
				func() error {
					return c.Job(ctx, filtered)
				},
				c.options.retryOptions...,
			)
			if err != nil {
				return err
			}

			return c.Broker.Commit(ctx, records)
		},
	)
}
