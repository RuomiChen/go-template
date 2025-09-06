package group_member

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type Service interface {
	JoinGroup(groupID, userID uint64, role string) error
	LeaveGroup(groupID, userID uint64) error
	GetGroupMembers(groupID uint64) ([]GroupMember, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}
func (s *service) WithTx(tx *gorm.DB) *service {
	return &service{
		repo: s.repo.WithTx(tx),
	}
}

// JoinGroup 加入群组
// role 可以是 "owner" | "admin" | "member"
func (s *service) JoinGroup(groupID, userID uint64, role string) error {
	exists, _, err := s.repo.IsMember(groupID, userID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("already a member")
	}

	if role == "" {
		role = "member"
	}

	member := &GroupMember{
		GroupID: groupID,
		UserID:  userID,
		Role:    role,
		Status:  1,
	}

	if err := s.repo.AddMember(member); err != nil {
		// 判断是否是外键约束失败
		if isForeignKeyError(err) {
			return errors.New("group does not exist")
		}
		return err
	}

	return nil
}

// 判断是否是 MySQL 外键约束错误
func isForeignKeyError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "1452")
}

// 退出群组
func (s *service) LeaveGroup(groupID, userID uint64) error {
	exists, member, err := s.repo.IsMember(groupID, userID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("not a member")
	}

	// 这里可以改为逻辑删除 Status = 0，也可以直接删除
	member.Status = 0
	return s.repo.RemoveMember(groupID, userID)
}

// 获取群组所有成员
func (s *service) GetGroupMembers(groupID uint64) ([]GroupMember, error) {
	return s.repo.GetMembersByGroup(groupID)
}
