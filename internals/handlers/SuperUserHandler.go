package handlers

import (
	"log"
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
		Username string `form:"username" binding:"required"`
		Email    string `form:"email" binding:"required,email"`
		Password string `form:"password" binding:"required,min=6"`
	}

	// Use ShouldBind to handle both form data and JSON
	if err := c.ShouldBind(&request); err != nil {
		errorMessage := "Invalid input data"
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			errorMessage = validationErr.Error()
		}
		log.Printf("Binding error: %v", err) // Log the actual error for debugging
		strategy.Respond(c, gin.H{"error": errorMessage}, http.StatusBadRequest)
		return
	}

	log.Printf("Received registration data: %+v", request)

	// Register superuser
	err := h.service.RegisterSuperuser(c.Request.Context(), request.Username, request.Email, request.Password)
	if err != nil {
		log.Printf("Service error: %v", err) // Log the actual error for debugging
		strategy.Respond(c, gin.H{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// For HTMX or HTML response, render "register_success.html"
	responseData := gin.H{"message": "Superuser registered successfully"}
	strategy.Respond(c, responseData, http.StatusOK)
}

func (h *SuperuserHandler) LoginSuperuser(c *gin.Context) {
	// Retrieve the response strategy from the context
	strategy := responses.GetResponseStrategy(c)

	var request struct {
		Email    string `form:"email" binding:"required,email"`
		Password string `form:"password" binding:"required"`
	}

	// Use ShouldBind to handle both form data and JSON
	if err := c.ShouldBind(&request); err != nil {
		errorMessage := "Invalid input data"
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			errorMessage = validationErr.Error()
		}
		log.Printf("Binding error: %v", err) // Log the actual error for debugging
		strategy.Respond(c, gin.H{"error": errorMessage}, http.StatusBadRequest)
		return
	}

	log.Printf("Received login data: %+v", request)

	// Authenticate superuser
	superuser, err := h.service.AuthenticateSuperuser(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		log.Printf("Authentication error: %v", err) // Log the actual error for debugging
		strategy.Respond(c, gin.H{"error": "Invalid email or password"}, http.StatusUnauthorized)
		return
	}

	log.Printf("Authentication successful for user: %+v", superuser)

	// Respond with the success message
	strategy.Respond(c, gin.H{"message": "Login successful", "user": superuser}, http.StatusOK)
}
