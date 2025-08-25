package admin

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindByUsername(username string) (*Admin, error)
	Create(admin *Admin) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByUsername(username string) (*Admin, error) {
	var admin Admin
	if err := r.db.Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *repository) Create(admin *Admin) error {
	return r.db.Create(admin).Error
}
