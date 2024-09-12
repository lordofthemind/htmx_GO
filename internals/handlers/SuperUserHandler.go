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
	// Get response strategy (JSON or HTML) based on request headers
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

		// Respond with error
		strategy.Respond(c, map[string]interface{}{
			"template": "register_error.html",
			"error":    errorMessage,
		}, http.StatusBadRequest)
		return
	}

	// Call the service to register the superuser
	err := h.service.RegisterSuperuser(c.Request.Context(), request.Username, request.Email, request.Password)
	if err != nil {
		strategy.Respond(c, map[string]interface{}{
			"template": "register_error.html",
			"error":    err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	// On success, respond with a success message
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

// 1. Profile Management
func (h *SuperuserHandler) ProfileViewHandler(c *gin.Context) {
	strategy := responses.GetResponseStrategy(c)
	strategy.Respond(c, map[string]interface{}{
		"template": "profile.html",
		"title":    "Profile",
		"user_id":  c.GetString("userID"), // Assuming userID is in context
	}, http.StatusOK)
}

func (h *SuperuserHandler) ProfileUpdateHandler(c *gin.Context) {
	strategy := responses.GetResponseStrategy(c)

	var request struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password"`
	}

	if err := c.ShouldBind(&request); err != nil {
		strategy.Respond(c, map[string]interface{}{
			"template": "profile_edit.html",
			"error":    "Invalid input data",
		}, http.StatusBadRequest)
		return
	}

	err := h.service.UpdateProfile(c.Request.Context(), request.Username, request.Password)
	if err != nil {
		strategy.Respond(c, map[string]interface{}{
			"template": "profile_edit.html",
			"error":    "Failed to update profile",
		}, http.StatusInternalServerError)
		return
	}

	strategy.Respond(c, map[string]interface{}{
		"template": "profile_success.html",
		"message":  "Profile updated successfully",
	}, http.StatusOK)
}

// 2. Password Reset
func (h *SuperuserHandler) PasswordResetRequestHandler(c *gin.Context) {
	strategy := responses.GetResponseStrategy(c)

	var request struct {
		Email string `form:"email" binding:"required,email"`
	}

	if err := c.ShouldBind(&request); err != nil {
		strategy.Respond(c, map[string]interface{}{
			"template": "password_reset_request.html",
			"error":    "Invalid email address",
		}, http.StatusBadRequest)
		return
	}

	// Generate reset link and send email
	err := h.service.SendPasswordResetEmail(c.Request.Context(), request.Email)
	if err != nil {
		strategy.Respond(c, map[string]interface{}{
			"template": "password_reset_request.html",
			"error":    "Failed to send reset email",
		}, http.StatusInternalServerError)
		return
	}

	strategy.Respond(c, map[string]interface{}{
		"template": "password_reset_sent.html",
		"message":  "Password reset email sent successfully",
	}, http.StatusOK)
}

func (h *SuperuserHandler) PasswordResetHandler(c *gin.Context) {
	strategy := responses.GetResponseStrategy(c)

	var request struct {
		Password string `form:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBind(&request); err != nil {
		strategy.Respond(c, map[string]interface{}{
			"template": "password_reset_form.html",
			"error":    "Invalid password",
		}, http.StatusBadRequest)
		return
	}

	err := h.service.ResetPassword(c.Request.Context(), c.Param("token"), request.Password)
	if err != nil {
		strategy.Respond(c, map[string]interface{}{
			"template": "password_reset_form.html",
			"error":    "Failed to reset password",
		}, http.StatusInternalServerError)
		return
	}

	strategy.Respond(c, map[string]interface{}{
		"template": "password_reset_success.html",
		"message":  "Password reset successful",
	}, http.StatusOK)
}

// 3. Two-Factor Authentication (2FA)
func (h *SuperuserHandler) Enable2FAHandler(c *gin.Context) {
	strategy := responses.GetResponseStrategy(c)

	// Render the 2FA enable page
	strategy.Respond(c, map[string]interface{}{
		"template": "2fa_enable.html",
		"title":    "Enable Two-Factor Authentication",
	}, http.StatusOK)
}

func (h *SuperuserHandler) Verify2FAHandler(c *gin.Context) {
	strategy := responses.GetResponseStrategy(c)

	var request struct {
		Code string `form:"code" binding:"required"`
	}

	if err := c.ShouldBind(&request); err != nil {
		strategy.Respond(c, map[string]interface{}{
			"template": "2fa_verify.html",
			"error":    "Invalid 2FA code",
		}, http.StatusBadRequest)
		return
	}

	// Verify the 2FA code
	err := h.service.Verify2FA(c.Request.Context(), c.GetString("userID"), request.Code)
	if err != nil {
		strategy.Respond(c, map[string]interface{}{
			"template": "2fa_verify.html",
			"error":    "Failed to verify 2FA code",
		}, http.StatusUnauthorized)
		return
	}

	strategy.Respond(c, map[string]interface{}{
		"template": "2fa_success.html",
		"message":  "2FA verified successfully",
	}, http.StatusOK)
}

// 4. User Roles and Permissions (RBAC)
func (h *SuperuserHandler) RoleManagementHandler(c *gin.Context) {
	strategy := responses.GetResponseStrategy(c)

	roles, err := h.service.ListRoles(c.Request.Context())
	if err != nil {
		strategy.Respond(c, map[string]interface{}{
			"template": "role_management.html",
			"error":    "Failed to retrieve roles",
		}, http.StatusInternalServerError)
		return
	}

	strategy.Respond(c, map[string]interface{}{
		"template": "role_management.html",
		"roles":    roles,
	}, http.StatusOK)
}

// 5. User Activity Logs
func (h *SuperuserHandler) UserActivityLogHandler(c *gin.Context) {
	strategy := responses.GetResponseStrategy(c)

	logs, err := h.service.GetUserActivityLogs(c.Request.Context(), c.GetString("userID"))
	if err != nil {
		strategy.Respond(c, map[string]interface{}{
			"template": "activity_log.html",
			"error":    "Failed to retrieve activity logs",
		}, http.StatusInternalServerError)
		return
	}

	strategy.Respond(c, map[string]interface{}{
		"template": "activity_log.html",
		"logs":     logs,
	}, http.StatusOK)
}

// 6. File Upload/Download
func (h *SuperuserHandler) FileUploadHandler(c *gin.Context) {
	strategy := responses.GetResponseStrategy(c)

	file, err := c.FormFile("file")
	if err != nil {
		strategy.Respond(c, map[string]interface{}{
			"template": "file_upload.html",
			"error":    "Failed to upload file",
		}, http.StatusBadRequest)
		return
	}

	// Save the file
	err = c.SaveUploadedFile(file, "./uploads/"+file.Filename)
	if err != nil {
		strategy.Respond(c, map[string]interface{}{
			"template": "file_upload.html",
			"error":    "Failed to save file",
		}, http.StatusInternalServerError)
		return
	}

	strategy.Respond(c, map[string]interface{}{
		"template": "file_upload_success.html",
		"message":  "File uploaded successfully",
	}, http.StatusOK)
}

func (h *SuperuserHandler) FileDownloadHandler(c *gin.Context) {
	strategy := responses.GetResponseStrategy(c)

	fileID := c.Param("file_id")
	filePath, err := h.service.GetFilePath(fileID)
	if err != nil {
		strategy.Respond(c, map[string]interface{}{
			"template": "file_download.html",
			"error":    "File not found",
		}, http.StatusNotFound)
		return
	}

	c.File(filePath)
}
