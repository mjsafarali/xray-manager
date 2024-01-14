package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"net/http"
	"shadowsocks-manager/internal/config"
	"shadowsocks-manager/internal/coordinator"
	"shadowsocks-manager/internal/database"
	"shadowsocks-manager/internal/http/handlers"
	"shadowsocks-manager/internal/http/handlers/v1"
	"shadowsocks-manager/internal/http/middleware"
	"shadowsocks-manager/internal/http/validator"
	"time"
)

type Server struct {
	Engine      *echo.Echo
	config      *config.Config
	logger      *zap.Logger
	coordinator *coordinator.Coordinator
	database    *database.Database
}

func New(config *config.Config, l *zap.Logger, c *coordinator.Coordinator, d *database.Database) *Server {
	e := echo.New()
	e.HideBanner = true
	e.Validator = validator.New()

	return &Server{Engine: e, config: config, logger: l, coordinator: c, database: d}
}

func (s *Server) Run() {
	s.Engine.Use(echoMiddleware.CORS())
	s.Engine.Use(middleware.Logger(s.logger))

	s.Engine.Static("/", "web")

	s.Engine.GET("/ssconf/*", handlers.SSConf(s.database))
	s.Engine.GET("/subscription/*", handlers.Subscription(s.database))
	s.Engine.GET("/profile", handlers.Profile())

	g1 := s.Engine.Group("/v1")
	g1.POST("/sign-in", v1.SignIn(s.database))

	g1.GET("/profile", v1.ProfileShow(s.database))
	g1.POST("/profile/reset", v1.ProfileReset(s.coordinator, s.database))

	g2 := s.Engine.Group("/v1")
	g2.Use(middleware.Authorize(s.database))

	g2.GET("/users", v1.UsersIndex(s.database))
	g2.POST("/users", v1.UsersStore(s.coordinator, s.database))
	g2.PUT("/users", v1.UsersUpdate(s.coordinator, s.database))
	g2.DELETE("/users/:id", v1.UsersDelete(s.coordinator, s.database))
	g2.PATCH("/users/:id/zero", v1.KeysZero(s.coordinator, s.database))

	g2.GET("/servers", v1.ServersIndex(s.database))
	g2.POST("/servers", v1.ServersStore(s.coordinator, s.database))
	g2.PUT("/servers", v1.ServersUpdate(s.coordinator, s.database))
	g2.DELETE("/servers/:id", v1.ServersDelete(s.coordinator, s.database))

	g2.GET("/settings", v1.SettingsShow(s.database))
	g2.POST("/settings", v1.SettingsUpdate(s.coordinator, s.database))
	g2.POST("/settings/import", v1.SettingsImport(s.coordinator, s.database))

	g2.GET("/stats", v1.StatsShow(s.database))
	g2.POST("/stats/servers/zero", v1.StatsZeroServers(s.database))
	g2.POST("/stats/users/zero", v1.StatsZeroUsers(s.database))
	g2.POST("/stats/users/delete-all", v1.StatsDeleteAllUsers(s.coordinator, s.database))

	go func() {
		address := fmt.Sprintf("%s:%d", s.config.HttpServer.Host, s.config.HttpServer.Port)
		if err := s.Engine.Start(address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatal("http server: failed to start", zap.String("address", address), zap.Error(err))
		}
	}()
}

func (s *Server) Shutdown() {
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Engine.Shutdown(c); err != nil {
		s.logger.Error("http server: failed to close", zap.Error(err))
	} else {
		s.logger.Debug("http server: closed successfully")
	}
}
