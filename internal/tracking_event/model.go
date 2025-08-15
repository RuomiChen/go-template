package tracking_event

import (
	"mvc/internal/common"
)

type TrackingEvent struct {
	common.BaseModel
	UserID string `gorm:"size:64;" json:"user_id"`
	From   string `gorm:"size:255" json:"from"`
	To     string `gorm:"size:255" json:"to"`
	Extra  string `gorm:"type:text" json:"extra"` // 存 JSON
	Action string `gorm:"size:50" json:"action"`
}
