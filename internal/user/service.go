package user

import (
	"errors"
	"mvc/internal/redis"
	"mvc/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
	existingUser, err := s.repo.GetByUsername(username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err // 查询异常直接返回错误
	}
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

func (s *Service) Login(c *fiber.Ctx, username, password string) (*LoginResponse, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil || user == nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 校验密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// // 记录
	// // 获取客户端信息
	// ip := c.IP()
	// ua := c.Get("User-Agent")
	// device := utils.parseDeviceFromUA(ua) // 你可以用 ua-parser-js 或 Go 的 UA 库实现

	// // 记录登录事件
	// extra, _ := json.Marshal(tracking_event.LoginExtra{
	// 	IP:        ip,
	// 	Device:    device,
	// 	UserAgent: ua,
	// 	Timestamp: time.Now(),
	// })
	// event := tracking_event.TrackingEvent{
	// 	UserID: u.ID,
	// 	From:   ip,
	// 	Action: "login",
	// 	Extra:  string(extra),
	// }
	// s.DB.Create(&event)

	token, err := utils.GenerateToken(user.ID, 0, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	_ = s.redisService.SaveKey(c.Context(), token, user.ID, time.Hour*3)

	// 返回结构体，包含token和user
	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"` // 或者具体 user 类型
}

func (s *Service) ChangePassword(userID, oldPassword, newPassword string) error {
	// 1. 查出用户
	user, err := s.repo.GetByID(userID)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}

	if oldPassword == newPassword {
		return errors.New("新密码不能和旧密码相同")
	}

	// 2. 校验旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("旧密码错误")
	}

	// 3. 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 4. 更新数据库
	updates := map[string]interface{}{
		"password": string(hashedPassword),
	}
	return s.repo.PartialUpdate(userID, updates)
}
