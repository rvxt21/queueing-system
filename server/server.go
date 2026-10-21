package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rvxt21/queueing-system/config"
	_ "github.com/rvxt21/queueing-system/docs" // registers the generated OpenAPI spec
	"github.com/rvxt21/queueing-system/server/handlers"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

const readHeaderTimeout = 5 * time.Second

type Server struct {
	config     *config.HTTP
	logger     *slog.Logger
	handlers   *handlers.Handler
	httpServer *http.Server
}

func NewServer(cfg *config.Config, logger *slog.Logger, h *handlers.Handler) *Server {
	return &Server{
		config:   &cfg.HTTP,
		logger:   logger,
		handlers: h,
	}
}

func (s *Server) Start(ctx context.Context) error {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/services", s.handlers.GetServices)
	router.POST("/tickets", s.handlers.RequestService)
	router.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))

	s.httpServer = &http.Server{
		Addr:              net.JoinHostPort(s.config.Host, s.config.Port),
		Handler:           router,
		ReadTimeout:       s.config.ReadTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	var lc net.ListenConfig

	ln, err := lc.Listen(ctx, "tcp", s.httpServer.Addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", s.httpServer.Addr, err)
	}

	go func() {
		if err := s.httpServer.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("http server stopped", "error", err)
		}
	}()

	s.logger.InfoContext(ctx, "http server started", "addr", ln.Addr().String())

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}

	return nil
}
