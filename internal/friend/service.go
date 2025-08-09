package friend

import (
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type FriendService struct {
	db *gorm.DB
}

func NewFriendService(db *gorm.DB) *FriendService {
	return &FriendService{db: db}
}

type AddFriendReq struct {
	UserID   string `json:"user_id"`
	FriendID string `json:"friend_id"`
}

func (s *FriendService) AddFriend(rawData json.RawMessage) error {
	var req AddFriendReq
	if err := json.Unmarshal(rawData, &req); err != nil {
		return errors.New("invalid add_friend data")
	}

	if req.UserID == "" || req.FriendID == "" {
		return errors.New("user_id and friend_id cannot be empty")
	}

	if req.UserID == req.FriendID {
		return errors.New("cannot add yourself as friend")
	}

	var count int64
	err := s.db.Model(&FriendRelation{}).
		Where("user_id = ? AND friend_id = ?", req.UserID, req.FriendID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("already friends")
	}

	relation := FriendRelation{
		UserID:   req.UserID,
		FriendID: req.FriendID,
	}
	if err := s.db.Create(&relation).Error; err != nil {
		return err
	}

	return nil
}
