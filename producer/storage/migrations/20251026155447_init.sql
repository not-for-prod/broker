-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE outbox
(
    id         BIGSERIAL PRIMARY KEY,
    topic      text  not null,
    partition  text  not null,
    headers    jsonb not null default '{}',
    payload    jsonb not null default '{}',
    created_at timestamp with time zone default now()
);

CREATE TABLE outbox_offset
(
    producer_name TEXT PRIMARY KEY,
    offset        BIGINT NOT NULL,
    updated_at    timestamp with time zone default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

drop table if exists outbox, outbox_offset cascade;
-- +goose StatementEnd
