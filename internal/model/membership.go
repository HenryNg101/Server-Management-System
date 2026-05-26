package model

type Membership struct {
	UserID   uint   `gorm:"not null"`
	User     User   `gorm:"constraint:OnDelete:CASCADE;"`
	ServerID uint   `gorm:"not null"`
	Server   Server `gorm:"constraint:OnDelete:CASCADE;"`
	UserRole string `gorm:"not null"`
}
