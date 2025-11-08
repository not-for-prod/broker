package kafka

import (
	"github.com/twmb/franz-go/pkg/kgo"
)

type Implementation struct {
	client *kgo.Client
}

func NewImplementation(client *kgo.Client) *Implementation {
	return &Implementation{
		client: client,
	}
}
