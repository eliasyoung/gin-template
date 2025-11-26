package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eliasyoung/gin-template/internal/cache"
	"github.com/eliasyoung/gin-template/internal/config"
	"github.com/eliasyoung/gin-template/internal/dal"
	"github.com/eliasyoung/gin-template/internal/logger"
	"github.com/eliasyoung/gin-template/internal/server"
	"go.uber.org/zap"
)

//	@title			Gin Template Api
//	@version		1.0
//	@description	This is the API service for Gin Template.
//	@host			localhost:7001
//	@BasePath		/

// @securityDefinitions.apikey	AdminAuthToken
// @in							header
// @name						Authorization
// @description				(For Admin Panel) Type "Bearer" followed by a space and the token obtained from the admin login endpoint.
func main() {
	cfg, err := config.Load()
	if err != nil {
		logger.Get().Fatal("err_config_init", zap.Error(err))
	}

	db := dal.Init(cfg.DBConfig)
	redisClient, err := cache.InitRedis(&cfg.CacheConfig)
	if err != nil {
		logger.Get().Fatal("err_cache_init", zap.Error(err))
	}

	srv := server.Setup(cfg, db, redisClient)

	go func() {
		srv.Run()
	}()

	// wait sigint sigterm signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Get().Info("Received shutdown signal...")

	// shutdown timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// gracefully shutdown gin server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Get().Fatal("Gin server shutdown failed", zap.Error(err))
	}

	logger.Get().Info("Application shut down gracefully.")

}
