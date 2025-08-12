package friend_relation

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type FriendRepo interface {
	Create(ctx context.Context, relation *FriendRelation) error
	Exists(ctx context.Context, userA, userB uint) (bool, error)
	Delete(ctx context.Context, userA, userB uint) error
	ListFriends(ctx context.Context, userID uint) ([]FriendRelation, error)
}

type friendRepo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) FriendRepo {
	return &friendRepo{db: db}
}

// 创建好友关系，保证 UserA < UserB
func (r *friendRepo) Create(ctx context.Context, relation *FriendRelation) error {
	if relation.UserA > relation.UserB {
		relation.UserA, relation.UserB = relation.UserB, relation.UserA
	}

	exist, err := r.Exists(ctx, relation.UserA, relation.UserB)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("friend relation already exists")
	}
	return r.db.WithContext(ctx).Create(relation).Error
}

func (r *friendRepo) Exists(ctx context.Context, userA, userB uint) (bool, error) {
	if userA > userB {
		userA, userB = userB, userA
	}
	var count int64
	err := r.db.WithContext(ctx).
		Model(&FriendRelation{}).
		Where("user_a = ? AND user_b = ?", userA, userB).
		Count(&count).Error
	return count > 0, err
}

func (r *friendRepo) Delete(ctx context.Context, userA, userB uint) error {
	if userA > userB {
		userA, userB = userB, userA
	}
	return r.db.WithContext(ctx).
		Where("user_a = ? AND user_b = ?", userA, userB).
		Delete(&FriendRelation{}).Error
}

// 列出 userID 所有好友关系（只查询包含 userID 的记录）
func (r *friendRepo) ListFriends(ctx context.Context, userID uint) ([]FriendRelation, error) {
	var relations []FriendRelation
	err := r.db.WithContext(ctx).
		Where("user_a = ? OR user_b = ?", userID, userID).
		Find(&relations).Error
	return relations, err
}
