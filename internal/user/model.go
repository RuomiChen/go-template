package user

import "mvc/internal/common"

type User struct {
	common.BaseModel
	Name  string `json:"name"`
	Email string `gorm:"type:varchar(50);unique"`
}
