package models

import "time"

type Post struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	User        User      `json:"-"`
}
