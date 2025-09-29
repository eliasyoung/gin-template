package server

import (
	"fmt"

	"github.com/eliasyoung/gin-template/internal/config"
	"github.com/eliasyoung/gin-template/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Server struct {
	router *gin.Engine
	port   string
}

func Setup(cfg *config.Config, db *gorm.DB, redisClient *redis.Client) *Server {

	router := gin.Default()

	exampleHandler := handler.InitExampleHandler(db, redisClient)
	router.GET("/ping", exampleHandler.HandleOnPing)

	return &Server{
		router: router,
		port:   cfg.ServerConfig.PORT,
	}
}

func (s *Server) Run() {
	s.router.Run(fmt.Sprintf(":%s", s.port))
}
