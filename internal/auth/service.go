package auth

import (
	"errors"
	"mvc/internal/redis"
	"mvc/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo         *Repository
	jwtSecret    string
	redisService redis.Service
}

func NewService(repo *Repository, jwtSecret string, redisService redis.Service) *Service {
	return &Service{repo: repo, jwtSecret: jwtSecret, redisService: redisService}
}

func (s *Service) Login(c *fiber.Ctx, username string, password string) (string, error) {
	auth, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", err
	}

	// 用 bcrypt 比对明文密码和数据库哈希密码
	if err := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(auth.ID, s.jwtSecret)
	if err != nil {
		return "", err
	}

	_ = s.redisService.SaveKey(c.Context(), token, auth.ID, time.Hour*3)

	return token, nil
}
func (s *Service) Register(c *fiber.Ctx, username string, password string) error {
	// 1. 检查用户名是否已存在
	_, err := s.repo.FindByUsername(username)
	if err == nil {
		return errors.New("username already exists")
	}

	// 2. bcrypt 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	auth := Auth{
		Username: username,
		Password: string(hashedPassword),
	}

	// 4. 调用 repo 新增用户
	if err := s.repo.Create(&auth); err != nil {
		return err
	}

	return nil
}
