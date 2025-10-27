package pgxv5

import (
	context "context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/not-for-prod/broker/models"
)

func (i *Implementation) ListRecords(ctx context.Context, limit uint64, offset uint64) ([]models.Event, error) {
	rows, err := i.pool.Query(
		ctx, `
		SELECT id, topic, partition, headers, payload, created_at
		FROM outbox
		ORDER BY id ASC
		LIMIT $1 OFFSET $2
	`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("query records: %w", err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var (
			id        uint64
			topic     string
			partition string
			headers   []byte
			payload   []byte
			createdAt time.Time
		)
		if err := rows.Scan(&id, &topic, &partition, &headers, &payload, &createdAt); err != nil {
			return nil, fmt.Errorf("scan record: %w", err)
		}

		var hdr map[string]string
		if err := json.Unmarshal(headers, &hdr); err != nil {
			hdr = map[string]string{}
		}

		events = append(
			events, models.Event{
				Topic:     topic,
				Partition: partition,
				Headers:   hdr,
				Body:      payload,
			},
		)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return events, nil
}
