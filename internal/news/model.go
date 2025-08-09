package news

type News struct {
	ID      uint `gorm:"primaryKey"`
	Content string
}
