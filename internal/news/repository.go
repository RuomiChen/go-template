package news

import (
	"mvc/internal/common"

	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]News, error)
	GetByID(id uint64) (*News, error)
	Create(news *News) error
	Update(news *News) error
	PartialUpdate(id uint64, updates map[string]interface{}) error
	Delete(id string) error
	GetPaged(page, pageSize int) ([]News, int64, error)
	GetTopNews(limit int) ([]News, error)

	GetByIDs(ids []string) ([]News, error)
	GetNewsByTag(tagID uint, limit, offset int) ([]News, error)
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

func (r *repository) GetByID(id uint64) (*News, error) {
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

func (r *repository) PartialUpdate(id uint64, updates map[string]interface{}) error {
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

// 根据 tagId 获取新闻列表
func (r *repository) GetNewsByTag(tagID uint, limit, offset int) ([]News, error) {
	var newsList []News
	err := r.db.Joins("JOIN news_tag nt ON nt.news_id = news.id").
		Where("nt.tag_id = ?", tagID).
		Preload("Tags"). // 加载新闻关联标签
		Limit(limit).Offset(offset).
		Find(&newsList).Error
	return newsList, err
}

func (r *repository) GetByIDs(ids []string) ([]News, error) {
	var list []News
	if len(ids) == 0 {
		return list, nil
	}
	err := r.db.Where("id IN ?", ids).Find(&list).Error
	return list, err
}
