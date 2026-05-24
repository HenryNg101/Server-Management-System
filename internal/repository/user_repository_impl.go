package repository

import (
	"github.com/HenryNg101/server-management-system/internal/model"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll() ([]model.User, error) {
	var result []model.User

	err := r.db.Model(&model.User{}).Find(&result).Error
	return result, err
}

func (r *userRepository) Create(user model.User) (model.User, error) {
	// user := models.User{Name: name}
	err := r.db.Create(&user).Error
	return user, err

	// err := r.db.QueryRow(
	// 	context.Background(),
	// 	"INSERT INTO users (name) VALUES ($1) RETURNING id",
	// 	user.Name,
	// ).Scan(&user.ID)

	// return user, err
}
