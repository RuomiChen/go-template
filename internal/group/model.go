package group

import (
	"mvc/internal/common"
	"os/user"
)

// 群组表
type Group struct {
	common.BaseModel
	Name    string `gorm:"size:100;not null;comment:群组名称"`
	OwnerID uint64 `gorm:"not null;index;comment:群主用户ID"`

	Owner user.User `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
