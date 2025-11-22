package migrations

import (
	"social-media-analyzer/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.Group{}, &models.Post{})
	return err
}
