package comment

import "gorm.io/gorm"

type Repository interface {
	Create(comment *Comment) error
	GetByPost(newsId uint64, limit, offset int) ([]Comment, error)
	GetReplies(parentID uint64) ([]Comment, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(comment *Comment) error {
	return r.db.Create(comment).Error
}

func (r *repository) GetByPost(newsId uint64, limit, offset int) ([]Comment, error) {
	var comments []Comment
	err := r.db.Where("news_id = ? AND parent_id IS NULL", newsId).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&comments).Error
	return comments, err
}

func (r *repository) GetReplies(parentID uint64) ([]Comment, error) {
	var comments []Comment
	err := r.db.Where("parent_id = ?", parentID).
		Order("created_at ASC").
		Find(&comments).Error
	return comments, err
}
