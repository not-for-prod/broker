package models

import (
	"context"
	"encoding/json"
)

type EventID string

func (id *EventID) String() string {
	return string(*id)
}

type Event struct {
	Ctx       context.Context
	ID        EventID
	Topic     string
	Partition string
	Headers   map[string]string
	Body      []byte
}

func (e *Event) Scan(a any) error {
	return json.Unmarshal(e.Body, a)
}
