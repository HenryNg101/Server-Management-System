package auth

type LoginRequest struct {
	Email    string `json:"email" example:"admin@example.com"`
	Password string `json:"password" example:"password123"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshResponse struct {
	NewRefreshToken string `json:"refresh_token"`
	NewAccessToken  string `json:"access_token"`
}

type RefreshData struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}
