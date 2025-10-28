package producer

import (
	"context"
	"time"

	"github.com/not-for-prod/broker"
)

type Producer struct {
	Broker    Broker
	Storage   Storage
	TxManager TxManager
	Logger    broker.Logger
	options   options
	stop      chan struct{}
}

func New(broker Broker, storage Storage, txManager TxManager, logger broker.Logger, opts ...Option) *Producer {
	options := defaultOptions
	for _, opt := range opts {
		opt.apply(&options)
	}

	return &Producer{
		Broker:    broker,
		Storage:   storage,
		TxManager: txManager,
		Logger:    logger,
		options:   options,
		stop:      make(chan struct{}),
	}
}

func (p *Producer) Start(ctx context.Context) error {
	ticker := time.NewTicker(p.options.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-p.stop:
			return nil
		case <-ticker.C:
			if err := p.produce(ctx); err != nil {
				p.Logger.Error("failed to produce events", "error", err.Error())
			}
		}
	}
}

func (p *Producer) Stop(_ context.Context) error {
	close(p.stop)
	return nil
}

func (p *Producer) Push(ctx context.Context, e ...broker.Event) error {
	return p.Storage.Push(ctx, e)
}

func (p *Producer) produce(ctx context.Context) error {
	return p.TxManager.Do(
		ctx, func(ctx context.Context) error {
			offset, err := p.Storage.GetOffset(ctx, p.options.producerName)
			if err != nil {
				return err
			}

			records, err := p.Storage.ListRecords(ctx, p.options.batchSize, offset)
			if err != nil {
				return err
			}

			err = p.Broker.Push(ctx, records)
			if err != nil {
				return err
			}

			err = p.Storage.CommitOffset(ctx, p.options.producerName, uint64(len(records)))
			if err != nil {
				return err
			}

			return nil
		},
	)
}
