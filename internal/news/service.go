package news

import (
	"fmt"
	"mvc/internal/news_like"
	"mvc/internal/redis"
	"mvc/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Service interface {
	GetNewsList(page, pageSize int) ([]News, int64, error)
	GetNewsDetail(id, userID uint64) (*News, error)
	AddNews(news *News) error
	DeleteNews(id string) error
	UpdateNews(id uint64, news *News) (*News, error)
	PartialUpdateNews(id uint64, updates map[string]interface{}) (*News, error)
	UploadImage(c *fiber.Ctx, saveDir string) (string, error)
	GetTopNews(limit int) ([]News, error)

	GetNewsByTag(tagID uint, limit, offset int) ([]News, error)
	GetNewsByIDs(ids []string) ([]News, error)
}

type service struct {
	repo            Repository
	hashStore       *utils.RedisHashStore
	newsLikeService *news_like.Service
}

func NewService(repo Repository, redisService redis.Service, newsLikeService *news_like.Service) Service {
	hashStore := utils.NewRedisHashStore(redisService, "imghash:", time.Hour*24*7)
	return &service{repo: repo, hashStore: hashStore, newsLikeService: newsLikeService}
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
func (s *service) GetNewsDetail(id, userID uint64) (*News, error) {
	news, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 点赞数量
	likeCount, _ := s.newsLikeService.CountLikes(news.ID)
	news.LikeCount = int(likeCount)
	// 当前用户是否点赞
	if userID != 0 {
		isLike, _ := s.newsLikeService.IsLiked(userID, news.ID)
		fmt.Print(isLike)
		news.IsLike = isLike
	}

	return news, nil
}
func (s *service) AddNews(news *News) error {
	return s.repo.Create(news)
}
func (s *service) DeleteNews(id string) error {
	return s.repo.Delete(id)
}
func (s *service) UpdateNews(id uint64, news *News) (*News, error) {
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

func (s *service) PartialUpdateNews(id uint64, updates map[string]interface{}) (*News, error) {
	if err := s.repo.PartialUpdate(id, updates); err != nil {
		return nil, err
	}
	return s.repo.GetByID(id)
}
func (s *service) UploadImage(c *fiber.Ctx, saveDir string) (string, error) {
	allowExts := []string{".jpg", ".jpeg", ".png"}
	return utils.UploadImageWithHashCheck(c, "image", saveDir, allowExts, s.hashStore)
}
func (s *service) GetTopNews(limit int) ([]News, error) {
	newsList, err := s.repo.GetTopNews(limit)
	if err != nil {
		return nil, err
	}
	return newsList, nil
}

func (s *service) GetNewsByTag(tagID uint, limit, offset int) ([]News, error) {
	return s.repo.GetNewsByTag(tagID, limit, offset)
}

func (s *service) GetNewsByIDs(ids []string) ([]News, error) {
	newsList, err := s.repo.GetByIDs(ids)
	if err != nil {
		return nil, err
	}
	return newsList, nil
}
