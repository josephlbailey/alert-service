create table alert
(
    id              integer generated always as identity primary key,
    external_id     uuid        not null,
    created_at      timestamptz not null,
    updated_at      timestamptz not null,
    message         text        not null,
    unique (external_id)
);