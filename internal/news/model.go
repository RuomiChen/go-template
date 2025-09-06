package news

import (
	"mvc/internal/common"
	"mvc/internal/tag"
)

type News struct {
	common.BaseModel
	Title   string    `json:"title" gorm:"type:varchar(255);not null;index"` // 标题
	Content string    `json:"content" gorm:"type:text;not null"`             // 内容
	Cover   string    `json:"cover,omitempty" gorm:"type:varchar(255)"`      // 封面图片
	Author  string    `json:"author,omitempty" gorm:"type:varchar(100)"`     // 作者/发布
	Views   uint      `json:"views" gorm:"default:0"`
	Tags    []tag.Tag `gorm:"many2many:news_tag;"` // 多对多关联 Tag
	// 虚拟字段，不映射数据库
	IsLike       bool `json:"is_like" gorm:"column:IsLike"`
	LikeCount    int  `json:"like_count" gorm:"column:LikeCount"`
	IsCollect    bool `json:"is_collect" gorm:"column:IsCollect"`
	CollectCount int  `json:"collect_count" gorm:"column:CollectCount"`
}
