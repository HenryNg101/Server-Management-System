package model

import "time"

type Agent struct {
	ID         uint   `gorm:"primaryKey"`
	ServerID   uint   `gorm:"not null;index"`
	APIKeyHash string `gorm:"not null;uniqueIndex"`
	CreatedAt  time.Time
}
