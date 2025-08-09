package auth

import (
	"errors"
	"mvc/internal/redis"
	"mvc/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Service struct {
	repo         *Repository
	jwtSecret    string
	redisService redis.Service
}

func NewService(repo *Repository, jwtSecret string, redisService redis.Service) *Service {
	return &Service{repo: repo, jwtSecret: jwtSecret, redisService: redisService}
}

func (s *Service) Login(c *fiber.Ctx, username, password string) (string, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil || user.Password != password { // 测试用：明文
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, s.jwtSecret)

	if err != nil {
		return "", err
	}
	_ = s.redisService.SaveToken(c.Context(), token, "123", time.Hour*3)

	return token, nil
}
