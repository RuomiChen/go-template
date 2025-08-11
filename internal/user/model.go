package user

import "mvc/internal/common"

type User struct {
	common.BaseModel
	Name  string
	Email string
}
