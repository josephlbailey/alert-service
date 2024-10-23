package db

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/josephlbailey/alert-service/internal/db/domain"
)

var (
	ErrAlertNotExists = errors.New("alert for the given external id not found")
)

type Store interface {
	domain.Querier
	CreateAlertTX(ctx context.Context, arg domain.CreateAlertParams) (*domain.Alert, error)
	UpdateAlertByIDTX(ctx context.Context, arg domain.UpdateAlertByIDParams) (*domain.Alert, error)
	DeleteAlertByIDTX(ctx context.Context, id int32) error
}

type AlertServiceStore struct {
	*domain.Queries
	db *pgxpool.Pool
}

func NewAlertServiceStore(db *pgxpool.Pool) Store {
	return &AlertServiceStore{
		db:      db,
		Queries: domain.New(db),
	}
}

func (store *AlertServiceStore) GetAlertByExternalID(ctx context.Context, externalID uuid.UUID) (*domain.Alert, error) {
	alert, err := store.Queries.GetAlertByExternalID(ctx, externalID)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, ErrAlertNotExists
		}

		return nil, err
	}

	if alert == nil {
		return nil, ErrAlertNotExists
	}

	return alert, nil
}

func (store *AlertServiceStore) CreateAlertTX(
	ctx context.Context,
	arg domain.CreateAlertParams,
) (*domain.Alert, error) {

	tx, err := store.db.Begin(context.Background())
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(context.Background())

	qtx := store.Queries.WithTx(tx)

	alert, err := qtx.CreateAlert(ctx, arg)

	if err != nil {
		return nil, err
	}
	if err = tx.Commit(context.Background()); err != nil {
		return nil, err
	}

	return alert, nil
}

func (store *AlertServiceStore) UpdateAlertByIDTX(
	ctx context.Context,
	arg domain.UpdateAlertByIDParams,
) (*domain.Alert, error) {

	tx, err := store.db.Begin(context.Background())
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(context.Background())

	qtx := store.Queries.WithTx(tx)

	alert, err := qtx.UpdateAlertByID(ctx, arg)

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(context.Background()); err != nil {
		return nil, err
	}

	return alert, nil
}

func (store *AlertServiceStore) DeleteAlertByIDTX(
	ctx context.Context,
	id int32,
) error {

	tx, err := store.db.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	qtx := store.Queries.WithTx(tx)

	err = qtx.DeleteAlertByID(ctx, id)

	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}
