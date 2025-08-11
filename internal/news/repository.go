package news

import "gorm.io/gorm"

type Repository interface {
	GetAll() ([]News, error)
	Create(user *News) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll() ([]News, error) {
	var newsList []News
	err := r.db.Find(&newsList).Error
	return newsList, err
}

func (r *repository) Create(news *News) error {
	return r.db.Create(news).Error
}
