package auth

type LoginRequest struct {
	Email    string `json:"email" example:"admin@example.com"`
	Password string `json:"password" example:"password123"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
