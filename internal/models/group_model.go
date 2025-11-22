package models

type Group struct {
	ID     uint   `gorm:"primaryKey"`
	Domain string `gorm:"type:text;not null"`
	Posts  []Post
}
