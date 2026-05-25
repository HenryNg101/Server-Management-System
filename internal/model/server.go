package model

import (
	"time"
)

type Server struct {
	ID          uint      `gorm:"primarykey"`
	Name        string    `gorm:"unique;not null"`
	Status      bool      `gorm:"not null"`                 // Server status at the time (Is it on or off)
	IPv4Address string    `gorm:"type:inet;not null"`       // IPv4 of the server
	Port        uint      `gorm:"check:ipv4_port <= 65535"` // Port of the server
	Protocol    string    `gorm:"not null"`                 // Network protocol that the server use, could be TCP, FTP, SSH, etc.
	CreatedAt   time.Time `gorm:"not null"`                 // Automatically managed by GORM for creation time
	LastUpdated time.Time `gorm:"not null"`                 // Last time the server got updated

	Memberships []Membership `gorm:"foreignKey:ServerID"` // Manage many-to-many relationship between servers and users
}
