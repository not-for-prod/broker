package nats

import "github.com/nats-io/nats.go"

type Implementation struct {
	js nats.JetStream
}

func NewImplementation(js nats.JetStream) *Implementation {
	return &Implementation{
		js: js,
	}
}
