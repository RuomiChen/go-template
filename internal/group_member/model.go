package group_member

import (
	"mvc/internal/common"
	"time"
)

// 群组成员表
type GroupMember struct {
	common.BaseModel
	GroupID    uint64     `gorm:"index;not null;comment:群组ID"`
	UserID     uint64     `gorm:"index;not null;comment:用户ID"`
	Role       string     `gorm:"type:enum('owner','admin','member');default:'member';comment:角色"`
	Alias      string     `gorm:"size:50;comment:群内昵称"`
	MutedUntil *time.Time `gorm:"comment:禁言截止时间"`
	Status     int        `gorm:"default:1;comment:1正常 0已踢出"`
}
