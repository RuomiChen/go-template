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
	FromUserID uint   `json:"from_user_id"`
	ToUserID   uint   `json:"to_user_id"`
	Message    string `json:"message,omitempty"`
}

func (s *FriendService) AddFriend(rawData json.RawMessage) error {
	var req AddFriendReq
	if err := json.Unmarshal(rawData, &req); err != nil {
		return errors.New("invalid add_friend data")
	}

	if req.FromUserID == 0 || req.ToUserID == 0 {
		return errors.New("user_id and friend_id cannot be empty")
	}

	if req.FromUserID == req.ToUserID {
		return errors.New("cannot add yourself as friend")
	}

	var count int64
	err := s.db.Model(&FriendRequest{}).
		Where("from_user_id = ? AND to_user_id = ?", req.FromUserID, req.ToUserID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("already friends")
	}

	relation := FriendRequest{
		Message:    req.Message,
		FromUserID: req.FromUserID,
		ToUserID:   req.ToUserID,
		Status:     RequestPending,
	}
	if err := s.db.Create(&relation).Error; err != nil {
		return err
	}

	return nil
}
