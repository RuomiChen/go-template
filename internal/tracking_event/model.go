package tracking_event

import (
	"encoding/json"
	"mvc/internal/common"
	"time"

	"gorm.io/gorm"
)

type Tracker struct {
	DB *gorm.DB
}
type TrackingEvent struct {
	common.BaseModel
	UserID string `gorm:"size:64;" json:"user_id"`
	From   string `gorm:"size:255" json:"from"`
	To     string `gorm:"size:255" json:"to"`
	Extra  string `gorm:"type:text" json:"extra"` // 存 JSON
	Action string `gorm:"size:50" json:"action"`
}
type LoginExtra struct {
	IP        string    `json:"ip"`
	Device    string    `json:"device"`
	UserAgent string    `json:"user_agent"`
	Timestamp time.Time `json:"timestamp"`
}

// NewTracker 返回 Tracker
func NewTracker(db *gorm.DB) *Tracker {
	return &Tracker{DB: db}
}

// 创建登录事件
func (t *Tracker) CreateLoginEvent(userID, ip string, extra LoginExtra) error {
	data, _ := json.Marshal(extra)
	event := TrackingEvent{
		UserID: userID,
		From:   ip,
		Action: "login",
		Extra:  string(data),
	}
	return t.DB.Create(&event).Error
}
