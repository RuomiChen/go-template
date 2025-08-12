package friend_relation

import (
	"context"
	"errors"
)

type FriendService struct {
	repo FriendRepo
}

func NewService(repo FriendRepo) *FriendService {
	return &FriendService{repo: repo}
}

func (s *FriendService) AddFriend(ctx context.Context, userA, userB uint) error {
	if userA == 0 || userB == 0 {
		return errors.New("user ids cannot be zero")
	}
	if userA == userB {
		return errors.New("cannot add yourself as friend")
	}
	return s.repo.Create(ctx, &FriendRelation{UserA: userA, UserB: userB})
}

func (s *FriendService) RemoveFriend(ctx context.Context, userA, userB uint) error {
	if userA == 0 || userB == 0 {
		return errors.New("user ids cannot be zero")
	}
	return s.repo.Delete(ctx, userA, userB)
}

func (s *FriendService) ListFriends(ctx context.Context, userID uint) ([]FriendRelation, error) {
	if userID == 0 {
		return nil, errors.New("user id cannot be zero")
	}
	return s.repo.ListFriends(ctx, userID)
}
