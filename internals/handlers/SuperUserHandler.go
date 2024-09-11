package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lordofthemind/htmx_GO/internals/responses"
	"github.com/lordofthemind/htmx_GO/internals/services"
)

type SuperuserHandler struct {
	service services.SuperuserService
}

func NewSuperuserHandler(service services.SuperuserService) *SuperuserHandler {
	return &SuperuserHandler{service: service}
}

func (h *SuperuserHandler) RegisterSuperuser(c *gin.Context) {
	// Get the response strategy from the context
	strategy := c.MustGet("responseStrategy").(responses.ResponseStrategy)

	var request struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	// Bind JSON request
	if err := c.ShouldBindJSON(&request); err != nil {
		strategy.Respond(c, gin.H{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Register superuser
	err := h.service.RegisterSuperuser(c.Request.Context(), request.Username, request.Email, request.Password)
	if err != nil {
		strategy.Respond(c, gin.H{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Prepare response data
	responseData := gin.H{"message": "Superuser registered successfully"}

	// Respond using the chosen strategy
	strategy.Respond(c, responseData, http.StatusOK)
}

func (h *SuperuserHandler) LoginSuperuser(c *gin.Context) {
	strategy := c.MustGet("responseStrategy").(responses.ResponseStrategy)

	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		strategy.Respond(c, gin.H{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	superuser, err := h.service.AuthenticateSuperuser(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		strategy.Respond(c, gin.H{"error": err.Error()}, http.StatusUnauthorized)
		return
	}

	strategy.Respond(c, gin.H{"message": "Login successful", "user": superuser}, http.StatusOK)
}
