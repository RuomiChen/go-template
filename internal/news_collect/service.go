package news_collect

import (
	"context"
	"mvc/internal/news"
	"time"
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
	like, err := s.repo.DeepFind(newsID, userID)
	if err != nil {
		return false, err
	}

	if like != nil && !like.DeletedAt.Valid {
		// 已点赞 -> 逻辑删除
		now := time.Now()
		like.DeletedAt.Time = now
		like.DeletedAt.Valid = true
		_ = s.repo.Update(like)
		// _ = s.statsRepo.DecrementLike(ctx, newsID)
		return false, nil
	}

	if like != nil && like.DeletedAt.Valid {
		// 撤销逻辑删除
		like.DeletedAt.Valid = false
		_ = s.repo.Update(like)
		// _ = s.statsRepo.IncrementLike(ctx, newsID)
		return true, nil
	}

	// 新增点赞
	_ = s.repo.Create(&NewsCollect{NewsID: newsID, UserID: userID})
	// _ = s.statsRepo.IncrementLike(ctx, newsID)
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
