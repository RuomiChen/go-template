package friend_relation

import "mvc/internal/common"

type FriendRelation struct {
	common.BaseModel
	UserA uint `gorm:"index;not null"` // keep smaller id in UserA to ensure uniqueness if you want
	UserB uint `gorm:"index;not null"`
}
