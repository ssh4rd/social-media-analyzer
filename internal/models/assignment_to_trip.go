package models

// import "gorm.io/gorm"

type AssignmentToTrip struct {
	ID             uint `gorm:"primaryKey"`
	MoneySpent     int `gorm:"not null"`
	EmployeeID     uint
	BusinessTripID uint
	Employee     Employee     `gorm:"foreignKey:EmployeeID;references:ID"`
	BusinessTrip BusinessTrip `gorm:"foreignKey:BusinessTripID;references:ID"`
}
