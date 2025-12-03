package models

import "time"

type Group struct {
	ID          uint      `gorm:"primaryKey"`
	Domain      string    `gorm:"type:text;not null;uniqueIndex"`
	Subscribers int       `gorm:"default:0"`
	ParsedAt    time.Time `gorm:"autoCreateTime:milli"`
	Posts       []Post
}
