package consumer

import (
	"context"
	"testing"

	"github.com/not-for-prod/broker"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	consumer Consumer
	inbox    map[uint64]bool
}

func (suite *TestSuite) SetupSuite() {
	suite.inbox = make(map[uint64]bool)
	suite.consumer = Consumer{
		Broker: &BrokerMock{
			CommitFunc: func(ctx context.Context, events []broker.Event) error {
				return nil
			},
			PullFunc: func(ctx context.Context, batchSize uint64) ([]broker.Event, error) {
				resp := make([]broker.Event, 0, batchSize)
				for i := uint64(0); i < batchSize; i++ {
					resp = append(
						resp, broker.Event{
							Ctx:       context.Background(),
							ID:        i,
							Topic:     "test",
							Partition: "test",
							Headers:   nil,
							Body:      []byte(`{"key":"value"}`),
						},
					)
				}
				return resp, nil
			},
		},
		Storage: &StorageMock{
			SetNXFunc: func(ctx context.Context, events []broker.Event) ([]bool, error) {
				resp := make([]bool, 0, len(events))
				for _, event := range events {
					_, ok := suite.inbox[event.ID]
					resp = append(resp, !ok)
					suite.inbox[event.ID] = true
				}
				return resp, nil
			},
		},
		TxManager: &TxManagerMock{DoFunc: func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}},
		Job: func(ctx context.Context, events []broker.Event) error {
			return nil
		},
		Logger: nil,
		options: options{
			batchSize:    5,
			interval:     0,
			retryOptions: nil,
		},
		stop: nil,
	}
}

func (suite *TestSuite) TestConsume() {
	err := suite.consumer.consume(context.Background())
	suite.Require().NoError(err)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
