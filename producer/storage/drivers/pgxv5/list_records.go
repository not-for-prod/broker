package pgxv5

import (
	context "context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/not-for-prod/broker"
	"go.opentelemetry.io/otel/propagation"
)

func (i *Implementation) ListRecords(ctx context.Context, limit uint64, offset uint64) ([]broker.Event, error) {
	rows, err := i.pool.Query(
		ctx, `
		SELECT id, topic, partition, headers, body, trace_carrier, created_at
		FROM outbox
		ORDER BY id ASC
		LIMIT $1 OFFSET $2
	`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("query records: %w", err)
	}
	defer rows.Close()

	var events []broker.Event
	for rows.Next() {
		var (
			id           uint64
			topic        string
			partition    string
			headers      []byte
			payload      []byte
			traceCarrier []byte
			createdAt    time.Time
		)
		if err := rows.Scan(&id, &topic, &partition, &headers, &payload, &traceCarrier, &createdAt); err != nil {
			return nil, fmt.Errorf("scan record: %w", err)
		}

		var hdr map[string]string
		if err := json.Unmarshal(headers, &hdr); err != nil {
			hdr = map[string]string{}
		}

		var mapCarrier propagation.MapCarrier
		if err := json.Unmarshal(payload, &mapCarrier); err != nil {
			mapCarrier = propagation.MapCarrier{}
		}

		events = append(
			events, broker.Event{
				Topic:     topic,
				Partition: partition,
				Headers:   hdr,
				Body:      payload,
				Ctx:       broker.ContextFromMapCarrier(mapCarrier),
			},
		)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return events, nil
}
