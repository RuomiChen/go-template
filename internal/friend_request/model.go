package friend_request

import "mvc/internal/common"

type FriendRequestStatus int

const (
	RequestPending   FriendRequestStatus = 0
	RequestAccepted  FriendRequestStatus = 1
	RequestRejected  FriendRequestStatus = 2
	RequestCancelled FriendRequestStatus = 3
)

type FriendRequest struct {
	common.BaseModel

	FromUserID uint                `gorm:"index;not null" json:"from_user_id"`
	ToUserID   uint                `gorm:"index;not null" json:"to_user_id"`
	Message    string              `gorm:"type:text" json:"message,omitempty"`
	Status     FriendRequestStatus `gorm:"type:tinyint;default:0" json:"status"`
}
