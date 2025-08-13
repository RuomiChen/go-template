package tag

import "mvc/internal/common"

type Tag struct {
	common.BaseModel
	Name string `json:"name" gorm:"type:varchar(255);not null;index"` // 标题
}
