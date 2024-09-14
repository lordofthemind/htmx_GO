package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/lordofthemind/htmx_GO/internals/configs"
	"github.com/lordofthemind/htmx_GO/internals/responses"
	"github.com/lordofthemind/htmx_GO/internals/services"
	"github.com/lordofthemind/htmx_GO/pkgs/tokens"
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

// Centralized error response handling
func (h *SuperuserHandler) handleError(c *gin.Context, template string, errorMessage string, statusCode int) {
	strategy := responses.GetResponseStrategy(c)
	strategy.Respond(c, map[string]interface{}{
		"template": template,
		"error":    errorMessage,
	}, statusCode)
}

// Centralized success response handling
func (h *SuperuserHandler) handleSuccess(c *gin.Context, template string, message string, statusCode int) {
	strategy := responses.GetResponseStrategy(c)
	strategy.Respond(c, map[string]interface{}{
		"template": template,
		"message":  message,
	}, statusCode)
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
		h.handleError(c, "register_error.html", errorMessage, http.StatusBadRequest)
		return
	}

	err := h.service.RegisterSuperuser(c.Request.Context(), request.Username, request.Email, request.Password)
	if err != nil {
		h.handleError(c, "register_error.html", err.Error(), http.StatusInternalServerError)
		return
	}

	h.handleSuccess(c, "register_success.html", "Superuser registered successfully", http.StatusOK)
}

func (h *SuperuserHandler) LoginSuperuserHandler(c *gin.Context) {
	var request struct {
		Email    string `form:"email" binding:"required,email"`
		Password string `form:"password" binding:"required"`
	}

	if err := c.ShouldBind(&request); err != nil {
		errorMessage := "Invalid input data"
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			errorMessage = validationErr.Error()
		}
		h.handleError(c, "login_error.html", errorMessage, http.StatusBadRequest)
		return
	}

	user, err := h.service.AuthenticateSuperuser(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		h.handleError(c, "login_error.html", "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := h.tokenManager.GenerateJWT(user.ID.String())
	if err != nil {
		h.handleError(c, "login_error.html", "Failed to generate token", http.StatusInternalServerError)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("SuperUserAuthorization", token, int(configs.TokenAccessDuration.Seconds()), "/", "", false, true)
	h.handleSuccess(c, "login_success.html", "Login successful", http.StatusOK)
}

func (h *SuperuserHandler) LogoutSuperuserHandler(c *gin.Context) {
	c.SetCookie("SuperUserAuthorization", "", -1, "/", "", false, true)
	h.handleSuccess(c, "index.html", "Logout successful", http.StatusOK)
}

func (h *SuperuserHandler) DashboardSuperuserHandler(c *gin.Context) {
	userID, _ := c.Get("userID") // Assuming user ID is in context from JWT middleware

	strategy := responses.GetResponseStrategy(c)
	strategy.Respond(c, map[string]interface{}{
		"template": "dashboard.html",
		"title":    "Dashboard",
		"user_id":  userID,
	}, http.StatusOK)
}

func (h *SuperuserHandler) TestTemplate(c *gin.Context) {
	strategy := responses.GetResponseStrategy(c)
	strategy.Respond(c, map[string]interface{}{"template": "test.html"}, http.StatusOK)
}

func (h *SuperuserHandler) ProfileViewHandler(c *gin.Context) {
	strategy := responses.GetResponseStrategy(c)
	strategy.Respond(c, map[string]interface{}{
		"template": "profile.html",
		"title":    "Profile",
		"user_id":  c.GetString("userID"),
	}, http.StatusOK)
}

func (h *SuperuserHandler) ProfileUpdateHandler(c *gin.Context) {
	var request struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password"`
	}

	// Bind the form data
	if err := c.ShouldBind(&request); err != nil {
		h.handleError(c, "profile_edit.html", "Invalid input data", http.StatusBadRequest)
		return
	}

	// Retrieve the user ID from the context
	userIDStr, exists := c.Get("userID")
	if !exists {
		h.handleError(c, "profile_edit.html", "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Convert user ID from string to uuid.UUID
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		h.handleError(c, "profile_edit.html", "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Call UpdateProfile with the correct arguments
	err = h.service.UpdateProfile(c.Request.Context(), userID, request.Username, request.Password)
	if err != nil {
		h.handleError(c, "profile_edit.html", "Failed to update profile", http.StatusInternalServerError)
		return
	}

	h.handleSuccess(c, "profile_success.html", "Profile updated successfully", http.StatusOK)
}

func (h *SuperuserHandler) PasswordResetRequestHandler(c *gin.Context) {
	var request struct {
		Email string `form:"email" binding:"required,email"`
	}

	if err := c.ShouldBind(&request); err != nil {
		h.handleError(c, "password_reset_request.html", "Invalid email address", http.StatusBadRequest)
		return
	}

	err := h.service.SendPasswordResetEmail(c.Request.Context(), request.Email)
	if err != nil {
		h.handleError(c, "password_reset_request.html", "Failed to send reset email", http.StatusInternalServerError)
		return
	}

	h.handleSuccess(c, "password_reset_sent.html", "Password reset email sent successfully", http.StatusOK)
}

func (h *SuperuserHandler) PasswordResetHandler(c *gin.Context) {
	var request struct {
		Password string `form:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBind(&request); err != nil {
		h.handleError(c, "password_reset_form.html", "Invalid password", http.StatusBadRequest)
		return
	}

	err := h.service.ResetPassword(c.Request.Context(), c.Param("token"), request.Password)
	if err != nil {
		h.handleError(c, "password_reset_form.html", "Failed to reset password", http.StatusInternalServerError)
		return
	}

	h.handleSuccess(c, "password_reset_success.html", "Password reset successful", http.StatusOK)
}

func (h *SuperuserHandler) Enable2FAHandler(c *gin.Context) {
	strategy := responses.GetResponseStrategy(c)
	strategy.Respond(c, map[string]interface{}{
		"template": "2fa_enable.html",
		"title":    "Enable Two-Factor Authentication",
	}, http.StatusOK)
}

// Verify2FAHandler verifies the 2FA code for the authenticated superuser.
func (h *SuperuserHandler) Verify2FAHandler(c *gin.Context) {
	var request struct {
		Code string `form:"code" binding:"required"`
	}

	// Bind the form data
	if err := c.ShouldBind(&request); err != nil {
		h.handleError(c, "2fa_verify.html", "Invalid 2FA code", http.StatusBadRequest)
		return
	}

	// Retrieve the user ID from the context
	userIDStr := c.GetString("userID")

	// Convert user ID from string to uuid.UUID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.handleError(c, "2fa_verify.html", "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// Call Verify2FA with the correct arguments
	err = h.service.Verify2FA(c.Request.Context(), userID, request.Code)
	if err != nil {
		h.handleError(c, "2fa_verify.html", "Failed to verify 2FA code", http.StatusUnauthorized)
		return
	}

	h.handleSuccess(c, "2fa_success.html", "2FA verified successfully", http.StatusOK)
}

func (h *SuperuserHandler) RoleManagementHandler(c *gin.Context) {
	roles, err := h.service.ListRoles(c.Request.Context())
	if err != nil {
		h.handleError(c, "role_management.html", "Failed to retrieve roles", http.StatusInternalServerError)
		return
	}

	strategy := responses.GetResponseStrategy(c)
	strategy.Respond(c, map[string]interface{}{
		"template": "role_management.html",
		"title":    "Role Management",
		"roles":    roles,
	}, http.StatusOK)
}

func (h *SuperuserHandler) UserActivityLogHandler(c *gin.Context) {
	strategy := responses.GetResponseStrategy(c)

	// Retrieve the user ID from the context and convert it to uuid.UUID
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		strategy.Respond(c, map[string]interface{}{
			"template": "activity_log.html",
			"error":    "Invalid user ID format",
		}, http.StatusBadRequest)
		return
	}

	// Fetch the user activity logs using the UUID
	logs, err := h.service.GetUserActivityLogs(c.Request.Context(), userID)
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
