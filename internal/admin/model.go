package admin

import "mvc/internal/common"

type Admin struct {
	common.BaseModel
	Username string `gorm:"type:varchar(50);unique" json:"username"`
	Password string `gorm:"type:varchar(100);" json:"password"`
	Name     string `gorm:"type:varchar(50);" json:"name"`
	Email    string `gorm:"type:varchar(50);"`
}
