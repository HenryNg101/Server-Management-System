package user

import (
	"github.com/HenryNg101/cron-scheduler/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll() ([]model.User, error) {
	var result []model.User

	err := r.db.Model(&model.User{}).Find(&result).Error
	return result, err
}
