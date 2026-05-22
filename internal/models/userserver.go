package models

type UserServer struct {
	UserID   uint `gorm:"not null"`
	User     User
	ServerID uint `gorm:"not null"`
	Server   Server
	UserRole string `gorm:"not null"`
}
