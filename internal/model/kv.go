package model

type KV struct {
	Key   string `gorm:"primaryKey;size:100;comment:配置键" json:"key"`
	Value string `gorm:"type:text;not null;comment:配置值" json:"value"`
}

func (KV) TableName() string {
	return "kv_store"
}
