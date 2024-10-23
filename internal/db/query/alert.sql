-- name: CreateAlert :one
insert into alert (
                     external_id,
                     created_at,
                     updated_at,
                     message
)
values ($1, $2, $3, $4)
returning *;

-- name: GetAlertByExternalID :one
select *
from alert
where external_id = $1;

-- name: UpdateAlertByID :one
update alert
set message = $1,
    updated_at = $2
where id = $3
returning *;

-- name: DeleteAlertByID :exec
delete from alert
where id = $1;
