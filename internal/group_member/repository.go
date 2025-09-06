package group_member

import (
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	AddMember(member *GroupMember) error
	RemoveMember(groupID, userID uint64) error
	GetMembersByGroup(groupID uint64) ([]GroupMember, error)
	IsMember(groupID, userID uint64) (bool, *GroupMember, error)

	WithTx(tx *gorm.DB) Repository
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) AddMember(member *GroupMember) error {
	return r.db.Create(member).Error
}

func (r *repository) RemoveMember(groupID, userID uint64) error {
	return r.db.Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&GroupMember{}).Error
}
func (r *repository) WithTx(tx *gorm.DB) Repository {
	return &repository{db: tx}
}
func (r *repository) GetMembersByGroup(groupID uint64) ([]GroupMember, error) {
	var members []GroupMember
	err := r.db.Where("group_id = ? AND status = 1", groupID).Preload("User").Find(&members).Error
	return members, err
}

func (r *repository) IsMember(groupID, userID uint64) (bool, *GroupMember, error) {
	var member GroupMember
	fmt.Print("groupID", groupID)
	fmt.Print("userID", userID)

	err := r.db.Where("group_id = ? AND user_id = ?", groupID, userID).First(&member).Error
	fmt.Print("member", member)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil, nil
		}
		return false, nil, err
	}
	return true, &member, nil
}
