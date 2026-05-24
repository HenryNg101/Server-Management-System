package model

import "time"

type User struct {
	UserID    uint      `gorm:"primarykey;not null"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}
