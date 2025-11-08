package sqlx

import (
	context "context"

	"github.com/not-for-prod/broker/producer/storage/drivers/sqlx/xo"
	"go.opentelemetry.io/otel"
)

func (i *Implementation) CommitOffset(ctx context.Context, producerName string, offset uint64) error {
	ctx, span := otel.Tracer("").Start(ctx, "outbox.commit_offset")
	defer span.End()

	xoOffset := &xo.OutboxOffset{
		ProducerName: producerName,
		Offset:       int64(offset),
	}

	err := xoOffset.Update(i.tr(ctx))
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}
