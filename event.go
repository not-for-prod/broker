package broker

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type Event struct {
	Ctx       context.Context
	ID        uint64
	Topic     string
	Partition string
	Headers   map[string]string
	Body      []byte
}

func NewEvent(
	ctx context.Context,
	topic string,
	partition string,
	headers map[string]string,
	body []byte,
) *Event {
	return &Event{
		Ctx:       ctx,
		ID:        0,
		Topic:     topic,
		Partition: partition,
		Headers:   headers,
		Body:      body,
	}
}

func (e *Event) MapCarrier() propagation.MapCarrier {
	mapCarrier := make(propagation.MapCarrier)
	otel.GetTextMapPropagator().Inject(e.Ctx, mapCarrier)

	return mapCarrier
}

func ContextFromMapCarrier(mapCarrier propagation.MapCarrier) context.Context {
	return otel.GetTextMapPropagator().Extract(context.Background(), mapCarrier)
}
