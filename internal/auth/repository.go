package auth

import "gorm.io/gorm"

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByUsername(username string) (*User, error) {
	var user User
	if err := r.db.Where("name = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
