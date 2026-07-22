package model

import "time"

type Agent struct {
	ID       uint `gorm:"column:id;primaryKey"`
	ServerID uint `gorm:"column:server_id;not null;uniqueIndex:idx_server_instance"`

	// Auth
	APIKey string `gorm:"column:api_key;not null;uniqueIndex"`

	// Identity of an agent's instance
	InstanceID string `gorm:"column:instance_id;not null;uniqueIndex:idx_server_instance"` // UUID from agent
	Hostname   string `gorm:"column:hostname;not null"`

	// Lifecycle
	Status     string     `gorm:"column:status;type:text;not null;default:'active'"` // active, revoked
	LastSeenAt *time.Time `gorm:"column:last_seen_at"`

	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime"`

	Server Server `gorm:"foreignKey:ServerID"`

	// Prevent duplicate agent instance per server
	// UNIQUE(server_id, instance_id)
}
