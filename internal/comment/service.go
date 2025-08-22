package comment

type Service interface {
	CreateComment(userID, NewsId uint64, parentID *uint64, content string) error
	GetComments(NewsId uint64, limit, offset int) ([]*CommentResponse, error)
}

type service struct {
	repo Repository
	// userService 可以在这里注入，用于查询用户信息
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateComment(userID, newsId uint64, parentID *uint64, content string) error {
	comment := &Comment{
		UserID:   userID,
		NewsId:   newsId,
		ParentID: parentID,
		Content:  content,
	}
	return s.repo.Create(comment)
}

func (s *service) GetComments(NewsId uint64, limit, offset int) ([]*CommentResponse, error) {
	comments, err := s.repo.GetByPost(NewsId, limit, offset)
	if err != nil {
		return nil, err
	}

	var result []*CommentResponse
	for _, c := range comments {
		// 组装子评论
		replies, _ := s.repo.GetReplies(c.ID)
		replyList := make([]*CommentResponse, 0, len(replies))
		for _, r := range replies {
			replyList = append(replyList, &CommentResponse{
				ID:        r.ID,
				NewsId:    r.NewsId,
				UserID:    r.UserID,
				ParentID:  r.ParentID,
				Content:   r.Content,
				LikeCount: r.LikeCount,
				CreatedAt: r.CreatedAt,
				// User: TODO 从 userService 查
			})
		}

		result = append(result, &CommentResponse{
			ID:         c.ID,
			NewsId:     c.NewsId,
			UserID:     c.UserID,
			ParentID:   c.ParentID,
			Content:    c.Content,
			LikeCount:  c.LikeCount,
			ReplyCount: c.ReplyCount,
			CreatedAt:  c.CreatedAt,
			Replies:    replyList,
			// User: TODO 从 userService 查
		})
	}
	return result, nil
}
