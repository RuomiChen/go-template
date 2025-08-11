package news

import "mvc/internal/common"

type News struct {
	common.BaseModel

	Content string
}
