package tag

import (
	"mvc/internal/redis"
	"mvc/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Service interface {
	GetTagList(page, pageSize int) ([]Tag, int64, error)
	GetTagDetail(id string) (*Tag, error)
	AddTag(tag *Tag) error
	DeleteTag(id string) error
	UpdateTag(id string, tag *Tag) (*Tag, error)
	PartialUpdateTag(id string, updates map[string]interface{}) (*Tag, error)
	UploadImage(c *fiber.Ctx, saveDir string) (string, error)
}

type service struct {
	repo      Repository
	hashStore *utils.RedisHashStore
}

func NewService(repo Repository, redisService redis.Service) Service {
	hashStore := utils.NewRedisHashStore(redisService, "imghash:", time.Hour*24*7)
	return &service{repo: repo, hashStore: hashStore}
}

func (s *service) GetTagList(page, pageSize int) ([]Tag, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return s.repo.GetPaged(page, pageSize)
}
func (s *service) GetTagDetail(id string) (*Tag, error) {
	return s.repo.GetByID(id)
}
func (s *service) AddTag(tag *Tag) error {
	return s.repo.Create(tag)
}
func (s *service) DeleteTag(id string) error {
	return s.repo.Delete(id)
}
func (s *service) UpdateTag(id string, tag *Tag) (*Tag, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	tag.ID = existing.ID
	if err := s.repo.Update(tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *service) PartialUpdateTag(id string, updates map[string]interface{}) (*Tag, error) {
	if err := s.repo.PartialUpdate(id, updates); err != nil {
		return nil, err
	}
	return s.repo.GetByID(id)
}
func (s *service) UploadImage(c *fiber.Ctx, saveDir string) (string, error) {
	allowExts := []string{".jpg", ".jpeg", ".png"}
	return utils.UploadImageWithHashCheck(c, "image", saveDir, allowExts, s.hashStore)
}
