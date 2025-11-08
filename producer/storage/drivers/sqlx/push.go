package sqlx

import (
	context "context"
	"encoding/json"
	"time"

	broker "github.com/not-for-prod/broker"
	"github.com/not-for-prod/broker/producer/storage/drivers/sqlx/xo"
	"go.opentelemetry.io/otel"
)

func (i *Implementation) Push(ctx context.Context, events []broker.Event) error {
	ctx, span := otel.Tracer("").Start(ctx, "outbox.push")
	defer span.End()

	xoEvents := make(xo.Outboxs, 0, len(events))

	for _, event := range events {
		traceCarrier, err := json.Marshal(event.MapCarrier())
		if err != nil {
			return err
		}

		headers, err := json.Marshal(event.Headers)
		if err != nil {
			return err
		}

		xoEvent := xo.Outbox{
			Topic:        event.Topic,
			Partition:    event.Partition,
			Headers:      headers,
			Body:         event.Body,
			TraceCarrier: traceCarrier,
			CreatedAt:    time.Now().UTC(),
		}

		xoEvents = append(xoEvents, xoEvent)
	}

	err := xoEvents.Insert(i.tr(ctx))
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}
