package user

import (
	"context"

	"github.com/HenryNg101/auth-service/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]model.User, error)
	Create(ctx context.Context, user model.User) (model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
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

func (r *userRepository) Create(ctx context.Context, user model.User) (model.User, error) {
	err := r.db.WithContext(ctx).Create(&user).Error
	return user, err
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var result model.User
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&result).
		Error
	return &result, err
}
