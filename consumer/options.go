package consumer

import (
	"time"

	"github.com/avast/retry-go"
)

const (
	defaultBatchSize = 500
	defaultInterval  = time.Second
)

type options struct {
	batchSize    uint64
	interval     time.Duration
	retryOptions []retry.Option
}

// Option overrides behavior of Connect.
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

var defaultOptions = options{
	batchSize:    defaultBatchSize,
	interval:     defaultInterval,
	retryOptions: []retry.Option{},
}

func WithBatchSize(batchSize uint64) Option {
	return optionFunc(
		func(o *options) {
			o.batchSize = batchSize
		},
	)
}

func WithInterval(interval time.Duration) Option {
	return optionFunc(
		func(o *options) {
			o.interval = interval
		},
	)
}

func WithRetryOptions(retryOptions ...retry.Option) Option {
	return optionFunc(
		func(o *options) {
			o.retryOptions = retryOptions
		},
	)
}
