-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

create table outbox
(
    id            bigserial primary key    not null,
    topic         text                     not null,
    partition     text                     not null,
    headers       jsonb                    not null default '{}',
    body          jsonb                    not null default '{}',
    trace_carrier jsonb                    not null default '{}',
    created_at    timestamp with time zone not null default now()
);

create table outbox_offset
(
    producer_name text                     not null,
    "offset"      bigint                   not null default 0,
    updated_at    timestamp with time zone not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

drop table if exists outbox;
-- +goose StatementEnd
