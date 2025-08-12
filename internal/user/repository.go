package user

import (
	"mvc/internal/common"

	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]User, error)
	GetByID(id string) (*User, error)
	GetByUsername(username string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	PartialUpdate(id string, updates map[string]interface{}) error
	Delete(id string) error
	GetPaged(page, pageSize int) ([]User, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetByUsername(username string) (*User, error) {
	var user User
	err := r.db.Model(&User{}).Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *repository) GetAll() ([]User, error) {
	var users []User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *repository) GetByID(id string) (*User, error) {
	var user User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *repository) Update(user *User) error {
	return r.db.Omit("created_at").Save(user).Error
}

func (r *repository) PartialUpdate(id string, updates map[string]interface{}) error {
	return r.db.Model(&User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Delete(&User{}, id).Error
}
func (r *repository) GetPaged(page, pageSize int) ([]User, int64, error) {
	return common.Paginate[User](r.db, page, pageSize)
}
