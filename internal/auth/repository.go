package auth

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByUsername(username string) (*Auth, error) {
	var auth Auth
	if err := r.db.Where("username = ?", username).First(&auth).Error; err != nil {
		return nil, err
	}
	return &auth, nil
}

func (r *Repository) Create(auth *Auth) error {
	return r.db.Create(auth).Error
}
