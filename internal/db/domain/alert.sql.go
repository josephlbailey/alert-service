// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: alert.sql

package domain

import (
	"context"
	"time"

	uuid "github.com/gofrs/uuid/v5"
)

const createAlert = `-- name: CreateAlert :one
insert into alert (
                     external_id,
                     created_at,
                     updated_at,
                     message
)
values ($1, $2, $3, $4)
returning id, external_id, created_at, updated_at, message
`

type CreateAlertParams struct {
	ExternalID uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Message    string
}

func (q *Queries) CreateAlert(ctx context.Context, arg CreateAlertParams) (*Alert, error) {
	row := q.db.QueryRow(ctx, createAlert,
		arg.ExternalID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Message,
	)
	var i Alert
	err := row.Scan(
		&i.ID,
		&i.ExternalID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Message,
	)
	return &i, err
}

const deleteAlertByID = `-- name: DeleteAlertByID :exec
delete from alert
where id = $1
`

func (q *Queries) DeleteAlertByID(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteAlertByID, id)
	return err
}

const getAlertByExternalID = `-- name: GetAlertByExternalID :one
select id, external_id, created_at, updated_at, message
from alert
where external_id = $1
`

func (q *Queries) GetAlertByExternalID(ctx context.Context, externalID uuid.UUID) (*Alert, error) {
	row := q.db.QueryRow(ctx, getAlertByExternalID, externalID)
	var i Alert
	err := row.Scan(
		&i.ID,
		&i.ExternalID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Message,
	)
	return &i, err
}

const updateAlertByID = `-- name: UpdateAlertByID :one
update alert
set message = $1,
    updated_at = $2
where id = $3
returning id, external_id, created_at, updated_at, message
`

type UpdateAlertByIDParams struct {
	Message   string
	UpdatedAt time.Time
	ID        int32
}

func (q *Queries) UpdateAlertByID(ctx context.Context, arg UpdateAlertByIDParams) (*Alert, error) {
	row := q.db.QueryRow(ctx, updateAlertByID, arg.Message, arg.UpdatedAt, arg.ID)
	var i Alert
	err := row.Scan(
		&i.ID,
		&i.ExternalID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Message,
	)
	return &i, err
}
