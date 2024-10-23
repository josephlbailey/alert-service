package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	cfg "github.com/josephlbailey/alert-service/config"
	"github.com/josephlbailey/alert-service/internal/api"
	"github.com/josephlbailey/alert-service/internal/db"
	l "github.com/josephlbailey/alert-service/internal/pkg/config"
)

func main() {

	var (
		config cfg.Config
		logger *zap.Logger
	)

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "dev"
		// set up logger for dev
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}

	defer logger.Sync()

	config = l.LoadConfig[cfg.Config]("alert-service", env)

	config.DB.Url = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.Database,
		config.DB.SslMode,
	)

	dbConn := db.Connect(config)
	defer db.Close(dbConn)

	db.AutoMigrate(config, logger)

	store := db.NewAlertServiceStore(dbConn)

	server := api.NewServer(
		config,
		logger,
		store,
	)

	server.MountHandlers()

	if config.Port == "" {
		config.Port = "8080"
	}

	addr := fmt.Sprintf(":%s", config.Port)

	// add graceful shutdown
	srv := &http.Server{
		Addr:    addr,
		Handler: server.Router(),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("error on listen...", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so no need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown Server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: ", zap.Error(err))
	}

	logger.Info("Server exiting")

}
