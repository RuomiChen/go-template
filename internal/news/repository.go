package news

import (
	"mvc/internal/common"

	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]News, error)
	GetByID(id string) (*News, error)
	Create(news *News) error
	Update(news *News) error
	PartialUpdate(id string, updates map[string]interface{}) error
	Delete(id string) error
	GetPaged(page, pageSize int) ([]News, int64, error)
	GetTopNews(limit int) ([]News, error)
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

func (r *repository) GetByID(id string) (*News, error) {
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
	return r.db.Omit("created_at").Save(news).Error
}

func (r *repository) PartialUpdate(id string, updates map[string]interface{}) error {
	return r.db.Model(&News{}).Where("id = ?", id).Updates(updates).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Delete(&News{}, id).Error
}
func (r *repository) GetPaged(page, pageSize int) ([]News, int64, error) {
	return common.Paginate[News](r.db, page, pageSize)
}
func (r *repository) GetTopNews(limit int) ([]News, error) {
	var newsList []News
	err := r.db.Order("views DESC").Limit(limit).Find(&newsList).Error
	if err != nil {
		return nil, err
	}
	return newsList, nil
}
