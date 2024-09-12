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

func (h *SuperuserHandler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func (h *SuperuserHandler) Register(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

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
		// Render an error template if needed
		strategy.Respond(c, gin.H{"template": "register_error.html", "error": errorMessage}, http.StatusBadRequest)
		return
	}

	err := h.service.RegisterSuperuser(c.Request.Context(), request.Username, request.Email, request.Password)
	if err != nil {
		// Render an error template if registration fails
		strategy.Respond(c, gin.H{"template": "register_error.html", "error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Success response with dynamic template
	strategy.Respond(c, gin.H{"template": "register_success.html", "message": "Superuser registered successfully"}, http.StatusOK)
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
		// Pass error template
		strategy.Respond(c, gin.H{"template": "login_error.html", "error": errorMessage}, http.StatusBadRequest)
		return
	}

	// Authenticate superuser
	superuser, err := h.service.AuthenticateSuperuser(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		// Pass error template
		strategy.Respond(c, gin.H{"template": "login_error.html", "error": "Invalid email or password"}, http.StatusUnauthorized)
		return
	}

	// Success response with dynamic template
	strategy.Respond(c, gin.H{"template": "login_response.html", "message": "Login successful", "user": superuser}, http.StatusOK)
}
