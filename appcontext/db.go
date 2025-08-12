package appcontext

import (
	"log"
	"mvc/internal/admin"
	"mvc/internal/friend_relation"
	"mvc/internal/friend_request"
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
		&admin.Admin{},
		&user.User{},
		&friend_request.FriendRequest{},
		&friend_relation.FriendRelation{},
		&news.News{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	return db
}
