package news_like

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Service interface {
	ToggleLike(ctx context.Context, newsID, userID uint64) (bool, error)
	IsLiked(userID, newsID uint64) (bool, error)
	CountLikes(newsID uint64) (int64, error)
}

type service struct {
	likeRepo Repository
	// 其他依赖
}

// , statsRepo NewsStatsRepo
func NewService(likeRepo Repository) Service {
	return &service{likeRepo: likeRepo}
	// , statsRepo: statsRepo
}
func (s *service) ToggleLike(ctx context.Context, newsID, userID uint64) (bool, error) {
	like, err := s.likeRepo.DeepFind(newsID, userID)
	if err != nil {
		return false, err
	}

	if like != nil && !like.DeletedAt.Valid {
		// 已点赞 -> 逻辑删除
		like.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}

		_ = s.likeRepo.Update(like)
		// _ = s.statsRepo.DecrementLike(ctx, newsID)
		return false, nil
	}

	if like != nil && like.DeletedAt.Valid {
		// 撤销逻辑删除
		like.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: false}

		_ = s.likeRepo.Update(like)
		return true, nil
	}

	// 新增点赞
	_ = s.likeRepo.Create(&NewsLike{NewsID: newsID, UserID: userID})
	// _ = s.statsRepo.IncrementLike(ctx, newsID)
	return true, nil
}
func (s *service) IsLiked(userID, newsID uint64) (bool, error) {
	return s.likeRepo.IsLiked(userID, newsID)
}

func (s *service) CountLikes(newsID uint64) (int64, error) {
	return s.likeRepo.CountLikes(newsID)
}
