package models

import (
	"time"

	// "gorm.io/gorm"
)

type BusinessTrip struct {
	ID          uint       `gorm:"primaryKey"`
	Destination string     `gorm:"type:text;not null"`
	StartAt     time.Time  `gorm:"type:date;not null"`
	EndAt       time.Time  `gorm:"type:date;not null"`
	Employees   []Employee `gorm:"many2many:assignment_to_trips;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
