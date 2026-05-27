package server

import (
	"context"
	"errors"
	"strings"

	"github.com/HenryNg101/server-management-system/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	FindAll(ctx context.Context, q GetServersQuery) ([]model.Server, int64, error)
	Create(ctx context.Context, server *model.Server) (*model.Server, error)
	FindByID(ctx context.Context, id uint, server *model.Server) (*model.Server, error)
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

	// Add filterings, one by one
	if q.Status != nil {
		db = db.Where("status = ?", *q.Status)
	}
	if q.Protocol != nil {
		db = db.Where("protocol = ?", *q.Protocol)
	}
	if q.Name != nil {
		db = db.Where("name ILIKE ?", "%"+*q.Name+"%") // Since it's just simple matching, just use this simple match
	}

	// Count BEFORE pagination, to have the true size of the results count
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting (basic safe whitelist)
	allowedSort := map[string]bool{
		"id":           true,
		"name":         true,
		"created_at":   true,
		"last_updated": true,
	}

	sortBy := "id"
	if allowedSort[q.SortBy] {
		sortBy = q.SortBy
	}

	order := "asc"
	if strings.ToLower(q.Order) == "desc" {
		order = "desc"
	}

	// Since attackers can inject SQL in this through the request's sortby field and order, it's better to have the above whitelist to filter all that
	db = db.Order(sortBy + " " + order)

	// 📄 Pagination
	offset := (q.Page - 1) * q.PageSize

	err := db.
		Limit(q.PageSize).
		Offset(offset).
		Find(&servers).Error

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

func (r *serverRepository) FindByID(ctx context.Context, id uint, server *model.Server) (*model.Server, error) {
	err := r.db.Model(&model.Server{}).
		Where("id = ?", id).
		Find(&server).
		Error
	return server, err
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
