package user

import (
	"time"

	"github.com/HenryNg101/users-service/internal/model"
)

type GetUsersResponse struct {
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Role      model.UserRole `json:"role"`
	CreatedAt time.Time      `json:"created_at"`
}
