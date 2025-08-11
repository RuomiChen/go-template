package news

import (
	"mvc/internal/common"

	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]News, error)
	GetByID(id uint) (*News, error)
	Create(news *News) error
	Update(news *News) error
	Delete(id uint) error
	GetPaged(page, pageSize int) ([]News, int64, error)
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

func (r *repository) GetByID(id uint) (*News, error) {
	var news News
	err := r.db.First(&news, id).Error
	if err != nil {
		return nil, err
	}
	return &news, nil
}

func (r *repository) Create(news *News) error {
	return r.db.Create(news).Error
}

func (r *repository) Update(news *News) error {
	return r.db.Save(news).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&News{}, id).Error
}
func (r *repository) GetPaged(page, pageSize int) ([]News, int64, error) {
	return common.Paginate[News](r.db, page, pageSize)
}
