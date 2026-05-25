package model

import "time"

type User struct {
	ID        uint      `gorm:"primarykey;not null"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"unique;not null"`
	CreatedAt time.Time `gorm:"not null"`

	Memberships []Membership `gorm:"foreignKey:UserID"` // Manage many-to-many relationship between servers and users
}
