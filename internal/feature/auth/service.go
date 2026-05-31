package auth

import (
	"context"
	"errors"

	"github.com/HenryNg101/server-management-system/internal/feature/user"
	"github.com/HenryNg101/server-management-system/internal/model"
)

type Service interface {
	ValidateUser(ctx context.Context, email string, password string) (*model.User, error)
}

type authService struct {
	userRepo user.Repository
}

func NewService(r user.Repository) Service {
	return &authService{userRepo: r}
}

func (s *authService) ValidateUser(ctx context.Context, email string, password string) (*model.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Plain text password for now (TODO: Must hash later)
	if user.Password != password {
		return nil, errors.New("invalid password")
	}

	return user, nil
}
