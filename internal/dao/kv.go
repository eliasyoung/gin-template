package dao

import (
	"github.com/eliasyoung/gin-template/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// KVDao KV存储数据访问接口
type KVDao interface {
	Get(tx *gorm.DB, key string) (string, error)
	Set(tx *gorm.DB, key, value string) error
	Delete(tx *gorm.DB, key string) error
	Exists(tx *gorm.DB, key string) (bool, error)
	GetAll(tx *gorm.DB) ([]model.KV, error)
	GetByPrefix(tx *gorm.DB, prefix string) ([]model.KV, error)
}

type GormKVDao struct{}

func NewGormKVDao() KVDao {
	return &GormKVDao{}
}

// Get 获取配置值
func (d *GormKVDao) Get(tx *gorm.DB, key string) (string, error) {
	var kv model.KV
	err := tx.Where("key = ?", key).First(&kv).Error
	if err != nil {
		return "", err
	}
	return kv.Value, nil
}

// Set 设置配置值
func (d *GormKVDao) Set(tx *gorm.DB, key, value string) error {
	kv := &model.KV{
		Key:   key,
		Value: value,
	}

	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(kv).Error
}

// Delete 删除配置
func (d *GormKVDao) Delete(tx *gorm.DB, key string) error {
	return tx.Where("key = ?", key).Delete(&model.KV{}).Error
}

// Exists 检查配置是否存在
func (d *GormKVDao) Exists(tx *gorm.DB, key string) (bool, error) {
	var count int64
	err := tx.Model(&model.KV{}).Where("key = ?", key).Count(&count).Error
	return count > 0, err
}

// GetAll 获取所有配置
func (d *GormKVDao) GetAll(tx *gorm.DB) ([]model.KV, error) {
	var kvs []model.KV
	err := tx.Find(&kvs).Error
	return kvs, err
}

// GetByPrefix 根据前缀获取配置
func (d *GormKVDao) GetByPrefix(tx *gorm.DB, prefix string) ([]model.KV, error) {
	var kvs []model.KV
	err := tx.Where("key LIKE ?", prefix+"%").Find(&kvs).Error
	return kvs, err
}
