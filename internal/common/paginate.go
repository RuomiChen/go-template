package common

import "gorm.io/gorm"

// Paginate 通用分页
func Paginate[T any](db *gorm.DB, page, pageSize int) (list []T, total int64, err error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 统计总数
	if err = db.Model(new(T)).Count(&total).Error; err != nil {
		return
	}

	offset := (page - 1) * pageSize
	err = db.Limit(pageSize).Offset(offset).Find(&list).Error
	return
}
