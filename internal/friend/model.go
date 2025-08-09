package friend

type FriendRelation struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	UserID   string `gorm:"index;not null"`
	FriendID string `gorm:"index;not null"`
}
