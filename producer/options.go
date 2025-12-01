package producer

import "time"

const (
	defaultProducerName = "default"
	defaultBatchSize    = 500
	defaultInterval     = time.Second
)

type options struct {
	producerName string
	batchSize    uint64
	interval     time.Duration
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
	producerName: defaultProducerName,
	batchSize:    defaultBatchSize,
	interval:     defaultInterval,
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

func WithProducerName(name string) Option {
	return optionFunc(
		func(o *options) {
			o.producerName = name
		},
	)
}
