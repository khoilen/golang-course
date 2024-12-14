package migrations

import (
	"log"
	"user-auth/models"

	"gorm.io/gorm"
)

func ResetMigrations(db gorm.DB) {
	var err error

	err = db.Migrator().DropTable(&models.User{}, &models.Post{}, &models.Like{}, &models.Comment{}, &models.Follow{})
	if err != nil {
		log.Fatalf("failed to drop tables: %v", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Post{}, &models.Like{}, &models.Comment{}, &models.Follow{})
	if err != nil {
		log.Fatalf("failed to migrate database schema: %v", err)
	}

	log.Println("Database schema has been reset and migrated successfully!")
}
