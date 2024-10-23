package api

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/josephlbailey/alert-service/config"
	"github.com/josephlbailey/alert-service/internal/db"
)

type Server struct {
	config   config.Config
	logger   *zap.Logger
	router   *gin.Engine
	store    db.Store
	accounts gin.Accounts
}

func NewServer(config config.Config, logger *zap.Logger, store db.Store) *Server {
	var engine *gin.Engine
	if config.Environment == "test" || config.Environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
		fmt.Printf("%v environment detected", config.Environment)
		engine = gin.New()
	} else {
		engine = gin.Default()
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

	var accounts = make(gin.Accounts)

	for _, user := range config.Users {
		accounts[user.Username] = user.Password
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:4200"}
	corsConfig.AllowHeaders = []string{"*"}
	engine.Use(cors.New(corsConfig))

	server := &Server{
		config:   config,
		logger:   logger,
		router:   engine,
		store:    store,
		accounts: accounts,
	}
	return server
}

func (s *Server) MountHandlers() {

	health := s.router.Group("/healthz")
	health.GET("", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"status": "ok",
		})
	})

	alert := s.router.Group("/alert")
	alert.POST("", gin.BasicAuth(s.accounts), s.CreateAlert)
	alert.GET("/:externalID", s.GetAlertByExternalID)
	alert.PUT("/:externalID", gin.BasicAuth(s.accounts), s.UpdateAlertByExternalID)
	alert.DELETE("/:externalID", gin.BasicAuth(s.accounts), s.DeleteAlertByExternalID)
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) Router() *gin.Engine {
	return s.router
}
