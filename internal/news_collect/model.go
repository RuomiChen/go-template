package news_collect

import "mvc/internal/common"

type NewsCollect struct {
	common.BaseModel
	NewsID uint64 `gorm:"index"`
	UserID uint64 `gorm:"index"`
}
