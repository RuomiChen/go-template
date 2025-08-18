package news_like

import "mvc/internal/common"

type NewsLike struct {
	common.BaseModel
	NewsID uint64 `gorm:"index"`
	UserID uint64 `gorm:"index"`
}
