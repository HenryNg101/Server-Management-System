package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/HenryNg101/auth-service/internal/model"
	"github.com/HenryNg101/auth-service/internal/user"

	// internalAuth "github.com/HenryNg101/server-management-system/internal/shared/auth"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	ValidateUser(ctx context.Context, email string, password string) (*model.User, error)
	GenerateRefreshToken(ctx context.Context, userID uint, role string) (string, error)
	StoreRefreshToken(ctx context.Context, userID uint, role string, refreshToken string, ttl time.Duration) error
	GetUserFromRefreshToken(ctx context.Context, refreshToken string) (*RefreshData, error)
	DeleteOldRefreshToken(ctx context.Context, key string) error
}

type authService struct {
	userRepo  user.Repository
	redisRepo RedisServerRepository
}

func NewService(r user.Repository, redis RedisServerRepository) Service {
	return &authService{userRepo: r, redisRepo: redis}
}

func (s *authService) ValidateUser(ctx context.Context, email string, password string) (*model.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Hash check
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (s *authService) GenerateRefreshToken(ctx context.Context, userID uint, role string) (string, error) {
	refreshToken, err := GenerateRefreshToken(userID, role)
	if err != nil {
		return "", err
	}
	err = s.StoreRefreshToken(ctx, userID, role, refreshToken, 7*24*time.Hour)
	return refreshToken, err
}

func (s *authService) StoreRefreshToken(ctx context.Context, userID uint, role string, refreshToken string, ttl time.Duration) error {
	value := RefreshData{
		UserID: userID,
		Role:   role,
	}
	bytes, _ := json.Marshal(value)
	key := fmt.Sprintf("refresh:%s", refreshToken)
	return s.redisRepo.StoreToken(ctx, key, bytes, ttl)
}

func (s *authService) GetUserFromRefreshToken(ctx context.Context, refreshToken string) (*RefreshData, error) {
	key := fmt.Sprintf("refresh:%s", refreshToken)
	val, err := s.redisRepo.GetUserInfo(ctx, key)
	if err != nil {
		return nil, err
	}

	var data RefreshData
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *authService) DeleteOldRefreshToken(ctx context.Context, refreshToken string) error {
	key := fmt.Sprintf("refresh:%s", refreshToken)
	return s.redisRepo.DeleteOldRefreshToken(ctx, key)
}
