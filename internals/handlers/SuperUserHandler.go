package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lordofthemind/htmx_GO/internals/services"
)

type SuperuserHandler struct {
	service services.SuperuserService
}

func NewSuperuserHandler(service services.SuperuserService) *SuperuserHandler {
	return &SuperuserHandler{service: service}
}

func (h *SuperuserHandler) RegisterSuperuser(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.RegisterSuperuser(c.Request.Context(), request.Username, request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Superuser registered successfully"})
}

func (h *SuperuserHandler) LoginSuperuser(c *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	superuser, err := h.service.AuthenticateSuperuser(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Implement session or JWT for maintaining login state (e.g., issue a JWT token)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": superuser})
}
