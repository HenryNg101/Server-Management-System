package auth

import (
	"net/http"
	"strings"

	"github.com/HenryNg101/server-management-system/internal/middleware/auth"
	internalAuth "github.com/HenryNg101/server-management-system/internal/middleware/auth"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Help people log in, and generate tokens for them
// Login godoc
// @Summary Login
// @Description Authenticate user and return JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 401 {object} map[string]string
// @Router /login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.ValidateUser(c, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := internalAuth.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	resp := LoginResponse{Token: token}
	c.JSON(http.StatusOK, &resp)
}

// TODO:
// Refresh godoc
// @Summary Refresh token
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} LoginResponse
// @Router /refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		c.JSON(401, gin.H{"error": "invalid token"})
		return
	}

	claims, err := auth.ParseToken(parts[1])
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid token"})
		return
	}

	newToken, err := auth.GenerateToken(claims.UserID, claims.Role)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(200, gin.H{"token": newToken})
}
