package pgxv5

import (
	context "context"
)

func (i *Implementation) CommitOffset(ctx context.Context, producerName string, offset int64) error {
	_, err := i.pool.Exec(
		ctx, `
		INSERT INTO outbox_offset (producer_name, "offset", updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (producer_name)
		DO UPDATE SET "offset" = EXCLUDED.offset, updated_at = NOW()
	`, producerName, offset,
	)
	if err != nil {
		return err
	}

	return nil
}
