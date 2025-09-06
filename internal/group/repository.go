package group

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(group *Group) error
	FindByID(id uint64) (*Group, error)
	FindAll() ([]Group, error)
	Update(group *Group) error
	Delete(id uint64) error

	Count(total *int64) error
	FindWithLimit(groups *[]Group, limit, offset int) error
	GetGroupsByUser(userID uint64) ([]Group, error)
	// 事务
	WithTx(tx *gorm.DB) Repository
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) WithTx(tx *gorm.DB) Repository {
	return &repository{db: tx}
}
func (r *repository) Create(group *Group) error {
	return r.db.Create(group).Error
}

func (r *repository) FindByID(id uint64) (*Group, error) {
	var group Group
	if err := r.db.Preload("Owner").First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *repository) FindAll() ([]Group, error) {
	var groups []Group
	if err := r.db.Preload("Owner").Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *repository) Update(group *Group) error {
	return r.db.Save(group).Error
}

func (r *repository) Delete(id uint64) error {
	return r.db.Delete(&Group{}, id).Error
}
func (r *repository) Count(total *int64) error {
	return r.db.Model(&Group{}).Count(total).Error
}

func (r *repository) FindWithLimit(groups *[]Group, limit, offset int) error {
	return r.db.Preload("Owner").Limit(limit).Offset(offset).Find(groups).Error
}

// 获取用户所在的所有群组
func (r *repository) GetGroupsByUser(userID uint64) ([]Group, error) {
	var groups []Group

	// 用户是群主
	err := r.db.
		Where("owner_id = ?", userID).
		Preload("Owner").
		Find(&groups).Error
	if err != nil {
		return nil, err
	}

	// 用户是普通成员或管理员
	var memberGroups []Group
	err = r.db.
		Joins("JOIN group_members gm ON gm.group_id = groups.id").
		Where("gm.user_id = ?", userID).
		Preload("Owner").
		Find(&memberGroups).Error
	if err != nil {
		return nil, err
	}

	// 合并结果，避免重复
	groupMap := make(map[uint64]Group)
	for _, g := range groups {
		groupMap[g.ID] = g
	}
	for _, g := range memberGroups {
		groupMap[g.ID] = g
	}

	result := make([]Group, 0, len(groupMap))
	for _, g := range groupMap {
		result = append(result, g)
	}
	return result, nil
}
