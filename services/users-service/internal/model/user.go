package model

import "time"

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

type User struct {
	ID        uint      `gorm:"primarykey;not null"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Role      UserRole  `gorm:"type:text;not null;check:role IN ('admin','user')"`
	CreatedAt time.Time `gorm:"not null"`

	Servers []Server `gorm:"foreignKey:UserID"`
}
