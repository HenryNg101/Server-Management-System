package model

import "time"

type Server struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint

	Name      string    `gorm:"not null;unique"`
	IPv4      string    `gorm:"column:ipv4_address;type:inet"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:last_updated;not null;autoUpdateTime"`

	User   User    `gorm:"foreignKey:UserID"`
	Agents []Agent `gorm:"foreignKey:ServerID"`

	// Optional: enforce unique per user
	// gorm:"uniqueIndex:idx_user_server_name"
}
