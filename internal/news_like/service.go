package news_like

import (
	"context"
	"time"
)

type Service struct {
	likeRepo Repository
	// statsRepo NewsStatsRepo
}

// , statsRepo NewsStatsRepo
func NewService(likeRepo Repository) *Service {
	return &Service{likeRepo: likeRepo}
	// , statsRepo: statsRepo
}
func (s *Service) ToggleLike(ctx context.Context, newsID, userID uint64) (bool, error) {
	like, err := s.likeRepo.DeepFind(newsID, userID)
	if err != nil {
		return false, err
	}

	if like != nil && !like.DeletedAt.Valid {
		// 已点赞 -> 逻辑删除
		now := time.Now()
		like.DeletedAt.Time = now
		like.DeletedAt.Valid = true
		_ = s.likeRepo.Update(like)
		// _ = s.statsRepo.DecrementLike(ctx, newsID)
		return false, nil
	}

	if like != nil && like.DeletedAt.Valid {
		// 撤销逻辑删除
		like.DeletedAt.Valid = false
		_ = s.likeRepo.Update(like)
		// _ = s.statsRepo.IncrementLike(ctx, newsID)
		return true, nil
	}

	// 新增点赞
	_ = s.likeRepo.Create(&NewsLike{NewsID: newsID, UserID: userID})
	// _ = s.statsRepo.IncrementLike(ctx, newsID)
	return true, nil
}
func (s *Service) IsLiked(userID, newsID uint64) (bool, error) {
	return s.likeRepo.IsLiked(userID, newsID)
}

func (s *Service) CountLikes(newsID uint64) (int64, error) {
	return s.likeRepo.CountLikes(newsID)
}
