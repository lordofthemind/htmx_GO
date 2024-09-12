package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/lordofthemind/htmx_GO/internals/responses"
	"github.com/lordofthemind/htmx_GO/internals/services"
)

type SuperuserHandler struct {
	service services.SuperuserService
}

func NewSuperuserHandler(service services.SuperuserService) *SuperuserHandler {
	return &SuperuserHandler{service: service}
}

// RegisterRoutes sets up the routes for the application.
func (h *SuperuserHandler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

// RegisterRoutes sets up the routes for the application.
func (h *SuperuserHandler) Register(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

// RegisterRoutes sets up the routes for the application.
func (h *SuperuserHandler) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func (h *SuperuserHandler) RegisterSuperuser(c *gin.Context) {
	// Retrieve the response strategy from the context
	strategy := responses.GetResponseStrategy(c)

	var request struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		errorMessage := "Invalid input data"
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			errorMessage = validationErr.Error()
		}
		strategy.Respond(c, gin.H{"error": errorMessage}, http.StatusBadRequest)
		return
	}

	err := h.service.RegisterSuperuser(c.Request.Context(), request.Username, request.Email, request.Password)
	if err != nil {
		strategy.Respond(c, gin.H{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	responseData := gin.H{"message": "Superuser registered successfully"}
	strategy.Respond(c, responseData, http.StatusOK)
}

func (h *SuperuserHandler) LoginSuperuser(c *gin.Context) {
	// Retrieve the response strategy from the context
	strategy := responses.GetResponseStrategy(c)

	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		errorMessage := "Invalid input data"
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			errorMessage = validationErr.Error()
		}
		strategy.Respond(c, gin.H{"error": errorMessage}, http.StatusBadRequest)
		return
	}

	// Authenticate superuser
	superuser, err := h.service.AuthenticateSuperuser(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		strategy.Respond(c, gin.H{"error": "Invalid email or password"}, http.StatusUnauthorized)
		return
	}

	strategy.Respond(c, gin.H{"message": "Login successful", "user": superuser}, http.StatusOK)
}
