package friend

import "mvc/internal/common"

type FriendRelation struct {
	common.BaseModel
	UserID   string `gorm:"index;not null" json:"user_id"`
	FriendID string `gorm:"index;not null" json:"friend_id"`
}
