package models

import "time"

type Post struct {
	ID          uint   `gorm:"primary_key"`
	Title       string `gorm:"not null"`
	Description string
	UserID      uint      `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
