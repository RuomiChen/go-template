package news_collect

import (
	"context"
	"mvc/internal/news"
	"time"

	"gorm.io/gorm"
)

type Service interface {
	ToggleCollect(ctx context.Context, newsID, userID uint64) (bool, error)
	IsCollected(userID, newsID uint64) (bool, error)
	CountCollects(newsID uint64) (int64, error)
	GetUserCollects(userID uint64) ([]*news.News, error)
}

type service struct {
	repo Repository
	// 其他依赖
	newsService news.Service
}

// , statsRepo NewsStatsRepo
func NewService(repo Repository, newsService news.Service) Service {
	return &service{repo: repo, newsService: newsService}
}
func (s *service) ToggleCollect(ctx context.Context, newsID, userID uint64) (bool, error) {
	collect, err := s.repo.DeepFind(newsID, userID)
	if err != nil {
		return false, err
	}

	if collect != nil && !collect.DeletedAt.Valid {
		// 已收藏 -> 逻辑删除
		collect.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
		if err := s.repo.Update(collect); err != nil {
			return false, err
		}
		return false, nil
	}

	if collect != nil && collect.DeletedAt.Valid {
		// 撤销逻辑删除 -> 直接清空
		collect.DeletedAt = gorm.DeletedAt{}
		if err := s.repo.Update(collect); err != nil {
			return false, err
		}
		return true, nil
	}

	// 新增收藏
	if err := s.repo.Create(&NewsCollect{NewsID: newsID, UserID: userID}); err != nil {
		return false, err
	}
	return true, nil
}
func (s *service) IsCollected(userID, newsID uint64) (bool, error) {
	return s.repo.IsLiked(userID, newsID)
}

func (s *service) CountCollects(newsID uint64) (int64, error) {
	return s.repo.CountLikes(newsID)
}
func (s *service) GetUserCollects(userID uint64) ([]*news.News, error) {
	collects, err := s.repo.FindByUser(userID)
	if err != nil {
		return nil, err
	}

	ids := make([]uint64, 0, len(collects))
	for _, c := range collects {
		ids = append(ids, c.NewsID)
	}

	return s.newsService.GetNewsByIDs(ids)
}
