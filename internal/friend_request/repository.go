package friend_request

import (
	"context"

	"gorm.io/gorm"
)

type friendRepo struct {
	db *gorm.DB
}

type FriendRepository interface {
	CountPendingRequest(ctx context.Context, fromUserID, toUserID uint) (int64, error)
	CreateFriendRequest(ctx context.Context, req *FriendRequest) error
}

func NewRepository(db *gorm.DB) FriendRepository {
	return &friendRepo{db: db}
}

func (r *friendRepo) CountPendingRequest(ctx context.Context, fromUserID, toUserID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&FriendRequest{}).
		Where("from_user_id = ? AND to_user_id = ?", fromUserID, toUserID).
		Count(&count).Error
	return count, err
}

func (r *friendRepo) CreateFriendRequest(ctx context.Context, req *FriendRequest) error {
	return r.db.WithContext(ctx).Create(req).Error
}
