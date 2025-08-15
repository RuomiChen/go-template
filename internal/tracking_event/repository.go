package tracking_event

import "gorm.io/gorm"

type Repository interface {
	addTrackingEvent(trackingEvent *TrackingEvent) error
	GetUserTrackingEvents(userID string, action string) ([]TrackingEvent, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) addTrackingEvent(trackingEvent *TrackingEvent) error {
	return r.db.Create(trackingEvent).Error
}

// 根据用户id获取纪录列表
func (r *repository) GetUserTrackingEvents(userID string, action string) ([]TrackingEvent, error) {
	var events []TrackingEvent

	// 子查询：每个 to 取最新一条，且 action 匹配
	subQuery := r.db.Model(&TrackingEvent{}).
		Select("MAX(id) as id").
		Where("user_id = ? AND action = ?", userID, action).
		Group("`to`") // 按新闻路径去重

	// 主查询：按 created_at 降序
	err := r.db.Model(&TrackingEvent{}).
		Where("id IN (?)", subQuery).
		Order("created_at DESC").
		Find(&events).Error

	return events, err
}
