package server

import (
	"context"
	"fmt"
	"net/http"

	_ "github.com/eliasyoung/gin-template/docs"
	"github.com/eliasyoung/gin-template/internal/config"
	"github.com/eliasyoung/gin-template/internal/handler"
	"github.com/eliasyoung/gin-template/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Server struct {
	httpServer *http.Server
}

func Setup(cfg *config.Config, db *gorm.DB, redisClient *redis.Client) *Server {

	gin.SetMode(cfg.ServerConfig.MODE)
	router := gin.Default()

	exampleHandler := handler.InitExampleHandler(db, redisClient)
	router.GET("/ping", exampleHandler.HandleOnPing)

	if gin.Mode() != gin.ReleaseMode {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServerConfig.PORT),
		Handler: router,
	}

	return &Server{
		httpServer: srv,
	}
}

func (s *Server) Run() {
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Get().Fatal("HTTP server ListenAndServe error", zap.Error(err))
	}
}

// Graceful Shutdown
func (s *Server) Shutdown(ctx context.Context) error {
	logger.Get().Info("Shuting down HTTP server...")
	return s.httpServer.Shutdown(ctx)
}
