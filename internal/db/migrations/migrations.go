package migrations

import (
	"TP_Andreev/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.Employee{}, &models.BusinessTrip{}, &models.AssignmentToTrip{})
	return err
}