package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/lordofthemind/htmx_GO/internals/configs"
	"github.com/lordofthemind/htmx_GO/internals/responses"
	"github.com/lordofthemind/htmx_GO/internals/services"
	"github.com/lordofthemind/htmx_GO/internals/tokens"
)

type SuperuserHandler struct {
	service      services.SuperuserService
	tokenManager tokens.TokenManager
}

func NewSuperuserHandler(service services.SuperuserService, tokenManager tokens.TokenManager) *SuperuserHandler {
	return &SuperuserHandler{
		service:      service,
		tokenManager: tokenManager,
	}
}

func (h *SuperuserHandler) IndexRender(c *gin.Context) {
	strategy := responses.GetResponseStrategy(c)
	strategy.Respond(c, map[string]interface{}{
		"template": "index.html",
		"title":    "Index",
	}, http.StatusOK)
}

func (h *SuperuserHandler) RegisterRender(c *gin.Context) {
	strategy := responses.GetResponseStrategy(c)
	strategy.Respond(c, map[string]interface{}{
		"template": "register.html",
		"title":    "Register",
	}, http.StatusOK)
}

func (h *SuperuserHandler) LoginRender(c *gin.Context) {
	strategy := responses.GetResponseStrategy(c)
	strategy.Respond(c, map[string]interface{}{
		"template": "login.html",
		"title":    "Login",
	}, http.StatusOK)
}

func (h *SuperuserHandler) RegisterSuperuserHandler(c *gin.Context) {
	strategy := responses.GetResponseStrategy(c)

	var request struct {
		Username string `form:"username" binding:"required"`
		Email    string `form:"email" binding:"required,email"`
		Password string `form:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBind(&request); err != nil {
		errorMessage := "Invalid input data"
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			errorMessage = validationErr.Error()
		}
		strategy.Respond(c, map[string]interface{}{
			"template": "register_error.html",
			"error":    errorMessage,
		}, http.StatusBadRequest)
		return
	}

	err := h.service.RegisterSuperuser(c.Request.Context(), request.Username, request.Email, request.Password)
	if err != nil {
		strategy.Respond(c, map[string]interface{}{
			"template": "register_error.html",
			"error":    err.Error(),
		}, http.StatusBadRequest)
		return
	}

	strategy.Respond(c, map[string]interface{}{
		"template": "register_success.html",
		"message":  "Superuser registered successfully",
	}, http.StatusOK)
}

// LoginSuperuser handles the login process for superusers.
func (h *SuperuserHandler) LoginSuperuserHandler(c *gin.Context) {
	strategy := responses.GetResponseStrategy(c)

	var request struct {
		Email    string `form:"email" binding:"required,email"`
		Password string `form:"password" binding:"required"`
	}

	if err := c.ShouldBind(&request); err != nil {
		errorMessage := "Invalid input data"
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			errorMessage = validationErr.Error()
		}
		strategy.Respond(c, map[string]interface{}{
			"template": "login_error.html",
			"error":    errorMessage,
		}, http.StatusBadRequest)
		return
	}

	// Check if the user exists
	user, err := h.service.AuthenticateSuperuser(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		strategy.Respond(c, map[string]interface{}{
			"template": "login_error.html",
			"error":    "Invalid email or password",
		}, http.StatusUnauthorized)
		return
	}

	// Generate JWT token using the TokenManager
	token, err := h.tokenManager.GenerateJWT(user.ID.Hex())
	if err != nil {
		strategy.Respond(c, map[string]interface{}{
			"template": "login_error.html",
			"error":    "Failed to generate token",
		}, http.StatusInternalServerError)
		return
	}

	// Set JWT as cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("SuperUserAuthorization", token, int(configs.TokenAccessDuration.Seconds()), "/", "", false, true)

	// Respond with success message
	strategy.Respond(c, map[string]interface{}{
		"template": "login_success.html",
		"message":  "Login successful",
		"user_id":  user.ID,
	}, http.StatusOK)
}

// Logout handler
func (h *SuperuserHandler) LogoutSuperuserHandler(c *gin.Context) {
	// Clear the JWT cookie
	c.SetCookie("SuperUserAuthorization", "", -1, "/", "", false, true)

	// Respond with success message
	strategy := responses.GetResponseStrategy(c)
	strategy.Respond(c, map[string]interface{}{
		"template": "index.html",
		"message":  "Logout successful",
	}, http.StatusOK)
}

// Dashboard renders the dashboard page using HTMX
func (h *SuperuserHandler) DashboardSuperuserHandler(c *gin.Context) {
	strategy := responses.GetResponseStrategy(c)

	// Pass user details or other required data to the dashboard
	userID, _ := c.Get("userID") // Assuming you set user ID in context in JWT middleware

	strategy.Respond(c, map[string]interface{}{
		"template": "dashboard.html",
		"title":    "Dashboard",
		"user_id":  userID,
	}, http.StatusOK)
}

func (h *SuperuserHandler) TestTemplate(c *gin.Context) {
	strategy := responses.GetResponseStrategy(c)

	strategy.Respond(c, map[string]interface{}{
		"template": "test.html",
	}, http.StatusOK)

}
