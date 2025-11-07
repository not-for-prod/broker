package sqlx

import (
	context "context"
	"encoding/json"
	"strconv"

	"github.com/Masterminds/squirrel"
	broker "github.com/not-for-prod/broker"
	"github.com/not-for-prod/broker/producer/storage/drivers/sqlx/xo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func (i *Implementation) ListRecords(ctx context.Context, limit uint64, offset uint64) ([]broker.Event, error) {
	ctx, span := otel.Tracer("").Start(ctx, "outbox.list_records")
	defer span.End()

	builder := sq.Select(
		xo.Outbox{}.
			SelectColumns()...,
	).
		From(xo.Table_Outbox_With_Alias).
		Where(squirrel.Gt{"outbox_id": offset}).
		Limit(limit)

	sqlstr, args, err := builder.ToSql()
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	var xoEvents []xo.Outbox

	err = i.tr(ctx).SelectContext(ctx, &xoEvents, sqlstr, args...)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	events := make([]broker.Event, 0, len(xoEvents))

	for _, xoEvent := range xoEvents {
		var mapCarrier propagation.MapCarrier

		err = json.Unmarshal(xoEvent.TraceCarrier, &mapCarrier)
		if err != nil {
			return nil, err
		}

		var headers map[string]string

		err = json.Unmarshal(xoEvent.Headers, &headers)
		if err != nil {
			return nil, err
		}

		event := broker.Event{
			Ctx:       broker.ContextFromMapCarrier(mapCarrier),
			ID:        broker.EventID(strconv.FormatInt(xoEvent.ID, 10)),
			Topic:     xoEvent.Topic,
			Partition: xoEvent.Partition,
			Headers:   headers,
			Body:      xoEvent.Body,
		}
		events = append(events, event)
	}

	return events, nil
}
