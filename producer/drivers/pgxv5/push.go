package pgxv5

import (
	context "context"
	"encoding/json"
	"fmt"

	"github.com/not-for-prod/broker/models"
)

func (s *Implementation) Push(ctx context.Context, events []models.Event) error {
	if len(events) == 0 {
		return nil
	}

	topics := make([]string, len(events))
	partitions := make([]string, len(events))
	headersArr := make([][]byte, len(events))
	payloads := make([][]byte, len(events))

	for i, e := range events {
		headersJSON, err := json.Marshal(e.Headers)
		if err != nil {
			return fmt.Errorf("marshal headers for event %d: %w", i, err)
		}
		topics[i] = e.Topic
		partitions[i] = e.Partition
		headersArr[i] = headersJSON
		payloads[i] = e.Body
	}

	_, err := s.pool.Exec(
		ctx, `
		INSERT INTO outbox (topic, partition, headers, payload)
		SELECT * FROM unnest($1::text[], $2::text[], $3::jsonb[], $4::jsonb[])
	`, topics, partitions, headersArr, payloads,
	)
	if err != nil {
		return fmt.Errorf("insert outbox batch with unnest: %w", err)
	}

	return nil
}
