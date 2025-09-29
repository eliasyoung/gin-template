package handler

import (
	"github.com/eliasyoung/gin-template/pkg"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ExampleHandler struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func InitExampleHandler(db *gorm.DB, redisClient *redis.Client) *ExampleHandler {
	return &ExampleHandler{
		db:          db,
		redisClient: redisClient,
	}
}

func (h *ExampleHandler) HandleOnPing(c *gin.Context) {
	c.JSON(200, pkg.SuccessResponse(gin.H{
		"message": "pong",
	}),
	)
}
