package news_like

import (
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	Create(like *NewsLike) error
	Delete(newsID, userID uint64) error
	Exists(newsID, userID uint64) (bool, error)
	Find(newsID, userID uint64) (*NewsLike, error)
	DeepFind(newsID, userID uint64) (*NewsLike, error)
	Update(like *NewsLike) error

	IsLiked(userID, newsID uint64) (bool, error)
	CountLikes(newsID uint64) (int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(like *NewsLike) error {
	return r.db.Create(like).Error
}
func (r *repository) Update(like *NewsLike) error {
	return r.db.Save(like).Error
}
func (r *repository) Delete(newsID, userID uint64) error {
	return r.db.
		Where("news_id = ? AND user_id = ?", newsID, userID).
		Delete(&NewsLike{}).Error
}
func (r *repository) DeepFind(newsID, userID uint64) (*NewsLike, error) {
	var like NewsLike
	err := r.db.Unscoped().
		Where("news_id = ? AND user_id = ?", newsID, userID).
		First(&like).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 没有找到记录
		}
		return nil, err
	}

	return &like, nil
}
func (r *repository) Find(newsID, userID uint64) (*NewsLike, error) {
	var like NewsLike
	err := r.db.
		Where("news_id = ? AND user_id = ?", newsID, userID).
		First(&like).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 没有找到记录
		}
		return nil, err
	}

	return &like, nil
}
func (r *repository) Exists(newsID, userID uint64) (bool, error) {
	var count int64
	err := r.db.
		Model(&NewsLike{}).
		Where("news_id = ? AND user_id = ?", newsID, userID).
		Count(&count).Error
	return count > 0, err
}
func (r *repository) IsLiked(userID, newsID uint64) (bool, error) {
	var count int64
	err := r.db.
		Model(&NewsLike{}).
		Where("user_id = ? AND news_id = ?", userID, newsID).
		Count(&count).Error
	fmt.Print("like count", count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repository) CountLikes(newsID uint64) (int64, error) {
	var count int64
	err := r.db.
		Model(&NewsLike{}).
		Where("news_id = ? AND deleted_at IS NULL", newsID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
