package friend_request

import (
	"context"
	"encoding/json"
	"errors"
)

type FriendService struct {
	repo FriendRepository
}

func NewService(repo FriendRepository) *FriendService {
	return &FriendService{repo: repo}
}

type AddFriendReq struct {
	FromUserID uint   `json:"from_user_id"`
	ToUserID   uint   `json:"to_user_id"`
	Message    string `json:"message,omitempty"`
}

func (s *FriendService) AddFriend(ctx context.Context, rawData json.RawMessage) error {
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

	count, err := s.repo.CountPendingRequest(ctx, req.FromUserID, req.ToUserID)
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
	if err := s.repo.CreateFriendRequest(ctx, &relation); err != nil {
		return err
	}

	return nil
}
