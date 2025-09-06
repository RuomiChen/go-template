package group

import (
	"mvc/internal/group_member"

	"gorm.io/gorm"
)

type Service interface {
	CreateGroupWithOwner(name string, ownerID uint64) (*Group, error)
	GetGroupList(page, pageSize int) ([]Group, int64, error)
	ListGroups() ([]Group, error)
	UpdateGroup(groupID uint64, name string) (*Group, error)
	DeleteGroup(id uint64) error
	GetGroupDetail(id uint64) (*Group, error)
	GetGroupsByUser(userID uint64) ([]Group, error)
}

type service struct {
	db   *gorm.DB
	repo Repository
}

func NewService(db *gorm.DB, repo Repository) Service {
	return &service{db: db, repo: repo}
}

// 创建群组并把创建者加入成员
func (s *service) CreateGroupWithOwner(name string, ownerID uint64) (*Group, error) {
	group := &Group{
		Name:    name,
		OwnerID: ownerID,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 先创建群组
		if err := tx.Create(group).Error; err != nil {
			return err
		}

		// 把群主加入成员表
		member := &group_member.GroupMember{
			GroupID: group.ID,
			UserID:  ownerID,
			Role:    "owner",
			Status:  1,
		}

		if err := tx.Create(member).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return group, nil
}
func (s *service) GetGroupsByUser(userID uint64) ([]Group, error) {
	return s.repo.GetGroupsByUser(userID)
}
func (s *service) GetGroupList(page, pageSize int) ([]Group, int64, error) {
	var groups []Group
	var total int64

	if err := s.repo.Count(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := s.repo.FindWithLimit(&groups, pageSize, offset); err != nil {
		return nil, 0, err
	}
	return groups, total, nil
}

func (s *service) ListGroups() ([]Group, error) {
	return s.repo.FindAll()
}
func (s *service) GetGroupDetail(id uint64) (*Group, error) {
	return s.repo.FindByID(id)
}
func (s *service) UpdateGroup(groupID uint64, name string) (*Group, error) {
	group, err := s.repo.FindByID(groupID)
	if err != nil {
		return nil, err
	}

	group.Name = name
	if err := s.repo.Update(group); err != nil {
		return nil, err
	}

	return group, nil
}

func (s *service) DeleteGroup(id uint64) error {
	return s.repo.Delete(id)
}
