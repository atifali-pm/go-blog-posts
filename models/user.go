package models

import "time"

type User struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	FirstName string    `gorm:"not null" json:"first_name"`
	LastName  string    `gorm:"not null" json:"last_name"`
	Avatar    string    `json:"avatar"`
	Phone     string    `gorm:"not null" json:"phone"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
