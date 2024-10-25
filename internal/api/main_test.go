package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	conf "github.com/josephlbailey/alert-service/config"
	"github.com/josephlbailey/alert-service/internal/db"
	common "github.com/josephlbailey/alert-service/internal/pkg/config"
)

func newTestServer(_ *testing.T, store db.Store) *Server {
	config := common.LoadConfig[conf.Config]("alert-service", "dev")
	logger := zap.NewNop()
	server := NewServer(config, logger, store)
	server.MountHandlers()
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
