package comment

import (
	"mvc/internal/common"
	"time"
)

type Comment struct {
	common.BaseModel
	NewsId     uint64  `gorm:"not null;index" json:"news_id"`    // 所属文章/视频/商品
	UserID     uint64  `gorm:"not null;index" json:"user_id"`    // 评论用户
	ParentID   *uint64 `gorm:"index" json:"parent_id,omitempty"` // 父级评论ID (null=一级评论)
	Content    string  `gorm:"type:text;not null" json:"content"`
	LikeCount  int64   `gorm:"default:0" json:"like_count"`  // 点赞数
	ReplyCount int64   `gorm:"default:0" json:"reply_count"` // 回复数（仅一级评论需要维护，方便快速查询）
}

// CommentResponse 用于 API 返回
type CommentResponse struct {
	ID         uint64    `json:"id"`
	NewsId     uint64    `json:"news_id"` // 所属文章/视频/商品
	UserID     uint64    `json:"user_id"`
	ParentID   *uint64   `json:"parent_id,omitempty"`
	Content    string    `json:"content"`
	LikeCount  int64     `json:"like_count"`
	ReplyCount int64     `json:"reply_count"`
	CreatedAt  time.Time `json:"created_at"`

	User    *UserDTO           `json:"user,omitempty"`    // 关联用户信息（不入库）
	Replies []*CommentResponse `json:"replies,omitempty"` // 子回复（不入库，接口组装）
}

// UserDTO 只是返回用的用户信息
type UserDTO struct {
	ID       uint64 `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}
