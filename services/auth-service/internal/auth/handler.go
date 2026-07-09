package auth

import (
	"fmt"
	"net/http"

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

	accessToken, err := GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to generate access token - %s\n", err.Error())})
		return
	}
	// Generate and store refresh token to Redis
	refreshToken, err := h.service.GenerateRefreshToken(c, user.ID, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to generate refresh token - %s\n", err.Error())})
		return
	}

	resp := LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}
	c.JSON(http.StatusOK, &resp)
}

// Refresh godoc
// @Summary Refresh token
// @Tags auth
// @Produce json
// @Param request body RefreshRequest true "Refresh token"
// @Success 200 {object} RefreshResponse
// @Router /refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	userInfo, err := h.service.GetUserFromRefreshToken(c, req.RefreshToken)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid refresh token"})
		return
	}

	// Delete old token (rotation) -> For security purposes
	// This is so that, when user refresh for a new token, potential attackers can't make use of the old one anymore
	h.service.DeleteOldRefreshToken(c, req.RefreshToken)

	// Issue new tokens
	accessToken, _ := GenerateAccessToken(userInfo.UserID, userInfo.Role)
	refreshToken, err := h.service.GenerateRefreshToken(c, userInfo.UserID, userInfo.Role)

	// Store new refresh token
	// err = h.service.StoreRefreshToken(c, userID, refreshToken, 7*24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, &RefreshResponse{
		NewAccessToken:  accessToken,
		NewRefreshToken: refreshToken,
	})
}

// Logout godoc
// @Summary Logout
// @Description Invalidate refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh token"
// @Success 200 {object} map[string]string
// @Router /logout [post]
func (h *Handler) Logout(c *gin.Context) {
	var req RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	err := h.service.DeleteOldRefreshToken(c, req.RefreshToken)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to logout"})
		return
	}

	c.JSON(200, gin.H{"message": "logged out"})
}
