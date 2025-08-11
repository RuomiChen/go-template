package news

import "mvc/internal/common"

type News struct {
	common.BaseModel
	Title   string `json:"title" gorm:"type:varchar(255);not null;index"` // 标题
	Content string `json:"content" gorm:"type:text;not null"`             // 内容
	Cover   string `json:"cover,omitempty" gorm:"type:varchar(255)"`      // 封面图片
	Author  string `json:"author,omitempty" gorm:"type:varchar(100)"`     // 作者/发布
	Views   uint   `json:"views" gorm:"default:0"`
}
