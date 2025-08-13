package tag

import "mvc/internal/common"

type Tag struct {
	common.BaseModel
	Title string `json:"title" gorm:"type:varchar(255);not null;index"` // 标题
}
