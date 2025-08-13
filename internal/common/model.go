package common

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 提供 创建时间 / 更新时间 / 逻辑删除
type BaseModel struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `json:"created_at"`     // 自动填充
	UpdatedAt time.Time      `json:"updated_at"`     // 自动更新
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 逻辑删除
}
