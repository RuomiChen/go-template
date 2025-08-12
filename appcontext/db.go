package appcontext

import (
	"log"
	"mvc/internal/auth"
	"mvc/internal/friend"
	"mvc/internal/news"
	"mvc/internal/user"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(
		&auth.Auth{},
		&user.User{},
		&friend.FriendRelation{},
		&news.News{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	return db
}
