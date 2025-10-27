package consumer

import (
	"time"

	"github.com/avast/retry-go"
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
	batchSize: 500,
	interval:  time.Second,
	retryOptions: []retry.Option{
		retry.Attempts(3),
		retry.Delay(time.Second),
	},
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
