package group

import (
	"mvc/internal/common"
	"mvc/internal/user"
)

// 群组类型
type GroupType string

const (
	GroupTypePublic  GroupType = "public"  // 公开群
	GroupTypePrivate GroupType = "private" // 私密群
)

// 群组表
type Group struct {
	common.BaseModel
	Name        string    `gorm:"size:100;not null;comment:群组名称" json:"name"`
	Description string    `gorm:"type:text;comment:群组描述" json:"description"`                                        // 新增群组描述
	Type        GroupType `gorm:"type:enum('public','private');default:'public';not null;comment:群组类型" json:"type"` // 新增群组类型
	OwnerID     uint64    `gorm:"not null;index;comment:群主用户ID" json:"owner_id"`

	Owner user.User `gorm:"foreignKey:OwnerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"owner"`
}
