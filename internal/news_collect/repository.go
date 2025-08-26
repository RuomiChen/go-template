package news_collect

import (
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	Create(collect *NewsCollect) error
	Delete(newsID, userID uint64) error
	Exists(newsID, userID uint64) (bool, error)
	Find(newsID, userID uint64) (*NewsCollect, error)
	DeepFind(newsID, userID uint64) (*NewsCollect, error)
	Update(collect *NewsCollect) error
	IsLiked(userID, newsID uint64) (bool, error)
	CountLikes(newsID uint64) (int64, error)

	// 👇 新增方法：查询某个用户的收藏列表
	FindByUser(userID uint64) ([]*NewsCollect, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(collect *NewsCollect) error {
	return r.db.Create(collect).Error
}
func (r *repository) Update(collect *NewsCollect) error {
	return r.db.Save(collect).Error
}
func (r *repository) Delete(newsID, userID uint64) error {
	return r.db.
		Where("news_id = ? AND user_id = ?", newsID, userID).
		Delete(&NewsCollect{}).Error
}
func (r *repository) DeepFind(newsID, userID uint64) (*NewsCollect, error) {
	var collect NewsCollect
	err := r.db.Unscoped().
		Where("news_id = ? AND user_id = ?", newsID, userID).
		First(&collect).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 没有找到记录
		}
		return nil, err
	}

	return &collect, nil
}
func (r *repository) Find(newsID, userID uint64) (*NewsCollect, error) {
	var collect NewsCollect
	err := r.db.
		Where("news_id = ? AND user_id = ?", newsID, userID).
		First(&collect).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 没有找到记录
		}
		return nil, err
	}

	return &collect, nil
}
func (r *repository) Exists(newsID, userID uint64) (bool, error) {
	var count int64
	err := r.db.
		Model(&NewsCollect{}).
		Where("news_id = ? AND user_id = ?", newsID, userID).
		Count(&count).Error
	return count > 0, err
}
func (r *repository) IsLiked(userID, newsID uint64) (bool, error) {
	var count int64
	err := r.db.
		Model(&NewsCollect{}).
		Where("user_id = ? AND news_id = ?", userID, newsID).
		Count(&count).Error
	fmt.Print("collect count", count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repository) CountLikes(newsID uint64) (int64, error) {
	var count int64
	err := r.db.
		Model(&NewsCollect{}).
		Where("news_id = ? AND deleted_at IS NULL", newsID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) FindByUser(userID uint64) ([]*NewsCollect, error) {
	var collects []*NewsCollect
	err := r.db.
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Find(&collects).Error
	if err != nil {
		return nil, err
	}
	return collects, nil
}
