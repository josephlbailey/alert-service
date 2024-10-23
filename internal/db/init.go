package db

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/josephlbailey/alert-service/config"
	"github.com/josephlbailey/alert-service/internal/pkg/path"
)

func Connect(config config.Config) *pgxpool.Pool {
	dbConfig, err := pgxpool.ParseConfig(config.DB.Url)
	if err != nil {
		log.Fatalf("Unable to parse database url: %v\n", config.DB.Url)
	}
	dbConfig.AfterConnect = func(ctx context.Context, pgconn *pgx.Conn) error {
		pgxuuid.Register(pgconn.TypeMap())
		return nil
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return pool
}

func Close(pool *pgxpool.Pool) {
	pool.Close()
}

func AutoMigrate(config config.Config, logger *zap.Logger) {

	mBase, err := path.Determine("internal")

	if err != nil {
		logger.Fatal("unable to determine migration path ", zap.Error(err))
	}

	mp := fmt.Sprintf("file://%s/%s", mBase, "internal/db/migration")
	dsn := fmt.Sprintf("pgx5://%s:%s@%s:%s/%s?sslmode=%s",
		config.DB.MigrationUsername,
		config.DB.MigrationPassword,
		config.DB.Host,
		config.DB.Port,
		config.DB.Database,
		config.DB.SslMode,
	)
	logger.Info("performing database migration...")
	m, err := migrate.New(mp, dsn)
	if err != nil {
		logger.Fatal("unable to create migration ", zap.Error(err))
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Fatal("unable to migrate database ", zap.Error(err))
	}
	logger.Info("database migration complete")
}
