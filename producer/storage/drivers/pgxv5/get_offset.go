package pgxv5

import (
	context "context"
	"database/sql"
	"errors"
)

func (i *Implementation) GetOffset(ctx context.Context, producerName string) (uint64, error) {
	var offset uint64

	err := i.pool.QueryRow(
		ctx, `
		SELECT "offset"
		FROM outbox_offset
		WHERE consumer_name = $1
	`, producerName,
	).Scan(&offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil // default offset
		}

		return 0, err
	}

	return offset, nil
}
