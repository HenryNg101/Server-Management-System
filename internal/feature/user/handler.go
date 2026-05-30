package user

import (
	"net/http"

	"github.com/HenryNg101/server-management-system/internal/model"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// CreateUser godoc
// @Summary Create a user
// @Description Create a new user, with specifications
// @Tags users
// @Param request body model.User true "User to be created"
// @Produce json
// @Success 200 {object} model.User
// @Router /users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.service.CreateUser(user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// GetUsers godoc
// @Summary Get all users
// @Description Retrieve list of users
// @Tags users
// @Produce json
// @Success 200 {array} model.User
// @Router /users [get]
func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.service.GetUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, users)
}
