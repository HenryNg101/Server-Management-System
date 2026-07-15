package model

import "time"

type Agent struct {
	ID        uint      `gorm:"primarykey;not null"`
	ServerID  uint      `gorm:"not null;index"`
	APIKey    string    `gorm:"column:api_key;not null;unique"`
	CreatedAt time.Time `gorm:"not null"`
	Server    Server
}
