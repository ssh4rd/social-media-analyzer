package models

type Post struct {
	ID        uint   `gorm:"primaryKey"`
	GroupID   uint
	Date      string `gorm:"type:text;not null"`
	Group     Group
	Views     int    `gorm:"not null"`
	Reactions int    `gorm:"not null"`
	Likes     int    `gorm:"not null"`
	Text      string `gorm:"type:text;not null"`
	Comments  int    `gorm:"not null"`
}
