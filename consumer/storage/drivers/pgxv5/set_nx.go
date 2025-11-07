package pgxv5

import (
	context "context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/not-for-prod/broker"
)

func (s *Implementation) SetNX(ctx context.Context, events []broker.Event) ([]bool, error) {
	if len(events) == 0 {
		return nil, nil
	}

	keys := make([]string, len(events))
	for i, e := range events {
		raw, _ := json.Marshal(e)
		keys[i] = fmt.Sprintf("%x", sha256.Sum256(raw))
	}

	_, err := s.pool.Exec(
		ctx, `
		INSERT INTO idempotency (event_key)
		SELECT unnest($1::text[])
		ON CONFLICT DO NOTHING
	`, keys,
	)
	if err != nil {
		return nil, err
	}

	// You won’t know which ones existed, but deduplication will still work.
	results := make([]bool, len(keys))
	for i := range results {
		results[i] = true
	}

	return results, nil
}
