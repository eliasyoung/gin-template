package main

import (
	"github.com/eliasyoung/gin-template/internal/cache"
	"github.com/eliasyoung/gin-template/internal/config"
	"github.com/eliasyoung/gin-template/internal/dal"
	"github.com/eliasyoung/gin-template/internal/logger"
	"github.com/eliasyoung/gin-template/internal/server"
	"go.uber.org/zap"
)

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
	srv.Run()
}
