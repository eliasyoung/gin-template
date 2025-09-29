package dal

import (
	"fmt"

	"github.com/eliasyoung/gin-template/internal/config"
	"github.com/eliasyoung/gin-template/internal/logger"
	"github.com/eliasyoung/gin-template/internal/model"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func Init(cfg config.DBConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	allModels := []interface{}{
		&model.KV{},
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger:      gormLogger.Default.LogMode(gormLogger.Info),
	})
	if err != nil {
		logger.Get().Panic("err_db_connection", zap.Error(err))
	}

	for _, m := range allModels {
		if !db.Migrator().HasTable(m) {
			if err = db.AutoMigrate(m); err != nil {
				logger.Get().Panic("err_db_migration", zap.Error(err))
			}
		}
		// if err = db.AutoMigrate(m); err != nil {
		// 	panic(err)
		// }
	}

	logger.Get().Info("Database connected")

	return db
}
