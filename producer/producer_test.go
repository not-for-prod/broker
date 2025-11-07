package producer

import (
	"context"
	"testing"

	"github.com/not-for-prod/broker"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	producer Producer
	offset   uint64
	buffer   []broker.Event
}

func (suite *TestSuite) SetupSuite() {
	suite.producer = Producer{
		Broker: &BrokerMock{
			PushFunc: func(ctx context.Context, r []broker.Event) error {
				return nil
			},
		},
		Storage: &StorageMock{
			CommitOffsetFunc: func(ctx context.Context, producerName string, offset uint64) error {
				suite.offset = offset
				return nil
			},
			GetOffsetFunc: func(ctx context.Context, producerName string) (uint64, error) {
				return suite.offset, nil
			},
			ListRecordsFunc: func(ctx context.Context, limit uint64, offset uint64) ([]broker.Event, error) {
				return suite.buffer, nil
			},
			PushFunc: func(ctx context.Context, events []broker.Event) error {
				for i, event := range events {
					event.ID = suite.offset + uint64(i) + 1
					suite.buffer = append(suite.buffer, event)
				}
				return nil
			},
		},
		TxManager: &TxManagerMock{DoFunc: func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}},
		Logger:  nil,
		options: defaultOptions,
		stop:    nil,
	}
}

func (suite *TestSuite) TestSend() {
	ctx := context.Background()

	err := suite.producer.Push(
		ctx, broker.Event{
			Ctx:       ctx,
			Topic:     "test",
			Partition: "test",
			Headers:   nil,
			Body:      []byte("test"),
		},
	)
	suite.Require().NoError(err)
	suite.Require().Len(suite.buffer, 1)

	err = suite.producer.produce(ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(suite.offset, uint64(len(suite.buffer)))
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
