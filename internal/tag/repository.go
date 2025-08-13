package tag

import (
	"mvc/internal/common"

	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]Tag, error)
	GetByID(id string) (*Tag, error)
	Create(tag *Tag) error
	Update(tag *Tag) error
	PartialUpdate(id string, updates map[string]interface{}) error
	Delete(id string) error
	GetPaged(page, pageSize int) ([]Tag, int64, error)
	GetTopTag(limit int) ([]Tag, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll() ([]Tag, error) {
	var tagList []Tag
	err := r.db.Find(&tagList).Error
	return tagList, err
}

func (r *repository) GetByID(id string) (*Tag, error) {
	var tag Tag
	err := r.db.First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *repository) Create(tag *Tag) error {
	return r.db.Create(tag).Error
}

func (r *repository) Update(tag *Tag) error {
	return r.db.Omit("created_at").Save(tag).Error
}

func (r *repository) PartialUpdate(id string, updates map[string]interface{}) error {
	return r.db.Model(&Tag{}).Where("id = ?", id).Updates(updates).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Delete(&Tag{}, id).Error
}
func (r *repository) GetPaged(page, pageSize int) ([]Tag, int64, error) {
	return common.Paginate[Tag](r.db, page, pageSize)
}
func (r *repository) GetTopTag(limit int) ([]Tag, error) {
	var tagList []Tag
	err := r.db.Order("views DESC").Limit(limit).Find(&tagList).Error
	if err != nil {
		return nil, err
	}
	return tagList, nil
}
