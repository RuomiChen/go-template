package friend

import "mvc/internal/common"

type FriendRelation struct {
	common.BaseModel
	UserID   string `gorm:"index;not null"`
	FriendID string `gorm:"index;not null"`
}
