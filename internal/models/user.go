package models

type User struct {
	UserID    uint   `gorm:"primarykey;not null"`
	Name      string `gorm:"not null"`
	Email     string `gorm:"not null"`
	CreatedAt string `gorm:"not null"`
}
