package user

import (
	"errors"
	"mvc/internal/redis"
	"mvc/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo         Repository
	jwtSecret    string
	redisService redis.Service
}

func NewService(repo Repository, jwtSecret string, redisService redis.Service) *Service {
	return &Service{repo: repo, jwtSecret: jwtSecret, redisService: redisService}
}

func (s *Service) GetUserList(page, pageSize int) ([]User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return s.repo.GetPaged(page, pageSize)
}
func (s *Service) GetUserDetail(id string) (*User, error) {
	return s.repo.GetByID(id)
}
func (s *Service) AddUser(user *User) error {
	return s.repo.Create(user)
}
func (s *Service) DeleteUser(id string) error {
	return s.repo.Delete(id)
}
func (s *Service) UpdateUser(id string, user *User) (*User, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	user.ID = existing.ID
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) PartialUpdateUser(id string, updates map[string]interface{}) (*User, error) {
	if err := s.repo.PartialUpdate(id, updates); err != nil {
		return nil, err
	}
	return s.repo.GetByID(id)
}
func (s *Service) Register(username, password string) error {
	// 先查用户名是否已存在
	existingUser, _ := s.repo.GetByUsername(username)
	if existingUser != nil {
		return errors.New("username already exists")
	}

	// 生成密码哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &User{
		Username: username,
		Password: string(hash),
		// 你还可以初始化其他字段
	}

	err = s.repo.Create(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Login(c *fiber.Ctx, username, password string) (string, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}
	if user == nil {
		return "", errors.New("invalid username or password")
	}

	// 校验密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	token, err := utils.GenerateToken(user.ID, 0, s.jwtSecret)
	if err != nil {
		return "", err
	}

	_ = s.redisService.SaveKey(c.Context(), token, user.ID, time.Hour*3)

	// 登录成功，返回 user 结构体
	return token, nil
}
