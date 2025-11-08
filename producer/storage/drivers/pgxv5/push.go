package pgxv5

import (
	context "context"
	"encoding/json"
	"fmt"

	"github.com/not-for-prod/broker"
)

func (s *Implementation) Push(ctx context.Context, events []broker.Event) error {
	if len(events) == 0 {
		return nil
	}

	topics := make([]string, len(events))
	partitions := make([]string, len(events))
	headersArr := make([][]byte, len(events))
	payloads := make([][]byte, len(events))
	traceCarriers := make([][]byte, len(events))

	for i, e := range events {
		headersJSON, err := json.Marshal(e.Headers)
		if err != nil {
			return fmt.Errorf("marshal headers for event %d: %w", i, err)
		}

		traceCarrier, err := json.Marshal(e.MapCarrier())
		if err != nil {
			return fmt.Errorf("marshal trace for event %d: %w", i, err)
		}

		topics[i] = e.Topic
		partitions[i] = e.Partition
		headersArr[i] = headersJSON
		payloads[i] = e.Body
		traceCarriers[i] = traceCarrier
	}

	_, err := s.pool.Exec(
		ctx, `
		INSERT INTO outbox (topic, partition, headers, body, trace_carrier)
		SELECT * FROM unnest($1::text[], $2::text[], $3::jsonb[], $4::jsonb[], $5::jsonb[])
	`, topics, partitions, headersArr, payloads, traceCarriers,
	)
	if err != nil {
		return fmt.Errorf("insert outbox batch with unnest: %w", err)
	}

	return nil
}
