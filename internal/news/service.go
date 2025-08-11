package news

import (
	"mvc/internal/redis"
	"mvc/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Service interface {
	GetNewsList(page, pageSize int) ([]News, int64, error)
	GetNewsDetail(id string) (*News, error)
	AddNews(news *News) error
	DeleteNews(id string) error
	UpdateNews(id string, news *News) (*News, error)
	PartialUpdateNews(id string, updates map[string]interface{}) (*News, error)
	UploadImage(c *fiber.Ctx) (string, error)
}

type service struct {
	repo      Repository
	hashStore *utils.RedisHashStore
}

func NewService(repo Repository, redisService redis.Service) Service {
	hashStore := utils.NewRedisHashStore(redisService, "imghash:", time.Hour*24*7)
	return &service{repo: repo, hashStore: hashStore}
}

func (s *service) GetNewsList(page, pageSize int) ([]News, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return s.repo.GetPaged(page, pageSize)
}
func (s *service) GetNewsDetail(id string) (*News, error) {
	return s.repo.GetByID(id)
}
func (s *service) AddNews(news *News) error {
	return s.repo.Create(news)
}
func (s *service) DeleteNews(id string) error {
	return s.repo.Delete(id)
}
func (s *service) UpdateNews(id string, news *News) (*News, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	news.ID = existing.ID
	if err := s.repo.Update(news); err != nil {
		return nil, err
	}
	return news, nil
}

func (s *service) PartialUpdateNews(id string, updates map[string]interface{}) (*News, error) {
	if err := s.repo.PartialUpdate(id, updates); err != nil {
		return nil, err
	}
	return s.repo.GetByID(id)
}
func (s *service) UploadImage(c *fiber.Ctx) (string, error) {
	allowExts := []string{".jpg", ".jpeg", ".png"}
	return utils.UploadImageWithHashCheck(c, "image", "./uploads", allowExts, s.hashStore)
}
