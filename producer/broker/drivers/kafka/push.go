package kafka

import (
	context "context"

	broker "github.com/not-for-prod/broker"
	"github.com/samber/lo"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.opentelemetry.io/otel"
)

func (i *Implementation) Push(ctx context.Context, events []broker.Event) error {
	ctx, span := otel.Tracer("").Start(ctx, "outbox.broker.push")
	defer span.End()

	records := make([]*kgo.Record, 0, len(events))

	for _, event := range events {
		record := convertEventToRecord(event)
		records = append(records, record)
	}

	produceResults := i.client.ProduceSync(ctx, records...)

	for _, res := range produceResults {
		if res.Err != nil {
			return res.Err
		}
	}

	return nil
}

func convertEventToRecord(e broker.Event) *kgo.Record {
	return &kgo.Record{
		Key:   []byte(e.Partition),
		Value: e.Body,
		Headers: lo.MapToSlice(
			e.Headers, func(key string, value string) kgo.RecordHeader {
				return kgo.RecordHeader{
					Key:   key,
					Value: []byte(value),
				}
			},
		),
		Topic:   e.Topic,
		Context: e.Ctx,
	}
}
