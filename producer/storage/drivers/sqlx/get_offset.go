package sqlx

import (
	context "context"

	"github.com/not-for-prod/broker/producer/storage/drivers/sqlx/xo"
	"go.opentelemetry.io/otel"
)

func (i *Implementation) GetOffset(ctx context.Context, producerName string) (uint64, error) {
	ctx, span := otel.Tracer("").Start(ctx, "outbox.get_offset")
	defer span.End()

	offset, err := xo.OutboxOffsetByProducerName(i.tr(ctx), producerName)
	if err != nil {
		span.RecordError(err)
		return 0, err
	}

	return uint64(offset.Offset), nil
}
