package user

import "gorm.io/gorm"

type Repository interface {
	GetAll() ([]User, error)
	Create(user *User) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll() ([]User, error) {
	var users []User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *repository) Create(user *User) error {
	return r.db.Create(user).Error
}
