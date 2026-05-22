package models

import (
	"net/netip"
	"time"
)

type Server struct {
	ServerID    uint         `gorm:"primarykey;not null"`  // Standard field for the primary key
	ServerName  string       `gorm:"unique;not null"`      // Foreign key to Users table
	Status      bool         `gorm:"type:string;not null"` // Server status at the time (Is it on or off)
	IPv4        netip.Prefix `gorm:"type:inet;not null"`   // IPv4 of the server
	CreatedAt   time.Time    `gorm:"type:time;not null"`   // Automatically managed by GORM for creation time
	LastUpdated time.Time    `gorm:"type:time;not null"`   // Last time the server got updated
}
