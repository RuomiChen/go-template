package auth

import (
	"errors"
	"mvc/pkg/utils"
)

type Service struct {
	repo      *Repository
	jwtSecret string
}

func NewService(repo *Repository, jwtSecret string) *Service {
	return &Service{repo: repo, jwtSecret: jwtSecret}
}

func (s *Service) Login(username, password string) (string, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil || user.Password != password { // 测试用：明文
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, s.jwtSecret)

	if err != nil {
		return "", err
	}

	return token, nil
}
