package server

import (
	"context"
	"errors"
	"strings"

	"github.com/HenryNg101/servers-service/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	FindAll(ctx context.Context, q GetServersQuery) ([]model.Server, int64, error)
	Create(ctx context.Context, server *model.Server) (*model.Server, error)
	FindByID(ctx context.Context, id uint) (*model.Server, error)
	Update(ctx context.Context, server *model.Server) (*model.Server, error)
	ExistsByID(ctx context.Context, id uint) (bool, error)
	Delete(ctx context.Context, id uint) error
}

type serverRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &serverRepository{db: db}
}

func (r *serverRepository) FindAll(ctx context.Context, q GetServersQuery) ([]model.Server, int64, error) {
	var servers []model.Server
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Server{})

	if q.UserID != nil {
		db = db.Where("user_id = ?", *q.UserID)
	}

	if q.Name != nil {
		db = db.Where("name ILIKE ?", "%"+*q.Name+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	allowedSort := map[string]bool{
		"id":         true,
		"name":       true,
		"created_at": true,
	}

	sortBy := "id"
	if allowedSort[q.SortBy] {
		sortBy = q.SortBy
	}

	order := "asc"
	if strings.ToLower(q.Order) == "desc" {
		order = "desc"
	}

	db = db.Order(sortBy + " " + order)

	if q.Page != nil && q.PageSize != nil {
		offset := (*q.Page - 1) * (*q.PageSize)
		db = db.Limit(*q.PageSize).Offset(offset)
	}

	err := db.Find(&servers).Error
	return servers, total, err
}

// Create if not exist, otherwise, returns error
func (r *serverRepository) Create(ctx context.Context, server *model.Server) (*model.Server, error) {
	result := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).
		Create(server)

	if result.Error != nil {
		return nil, result.Error
	}
	// Handle a failure to insert due to record already exist, not validations or something
	if result.RowsAffected == 0 {
		return nil, errors.New("The server with this name is already existed")
	}
	return server, nil
}

func (r *serverRepository) FindByID(ctx context.Context, id uint) (*model.Server, error) {
	var server model.Server

	err := r.db.WithContext(ctx).
		First(&server, id).Error

	return &server, err
}

func (r *serverRepository) Update(ctx context.Context, server *model.Server) (*model.Server, error) {
	err := r.db.WithContext(ctx).
		Save(server).Error

	return server, err
}

func (r *serverRepository) ExistsByID(ctx context.Context, id uint) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&model.Server{}).
		Where("id = ?", id).
		Count(&count).Error

	return count > 0, err
}

func (r *serverRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Delete(&model.Server{}, id).Error
}
