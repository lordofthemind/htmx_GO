package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lordofthemind/htmx_GO/internals/handlers"
	"github.com/lordofthemind/htmx_GO/pkgs/middlewares"
	"github.com/lordofthemind/htmx_GO/pkgs/tokens"
)

func RegisterSuperuserRoutes(router *gin.Engine, superuserHandler *handlers.SuperuserHandler, tokenManager tokens.TokenManager) {
	// Group for superuser-related routes
	superuserRoutes := router.Group("/superuser")
	{
		// Public routes
		superuserRoutes.GET("/", superuserHandler.IndexRender)
		superuserRoutes.GET("/register", superuserHandler.RegisterRender)
		superuserRoutes.GET("/login", superuserHandler.LoginRender)
		superuserRoutes.POST("/register", superuserHandler.RegisterSuperuserHandler)
		superuserRoutes.POST("/login", superuserHandler.LoginSuperuserHandler)

		// Apply JWTAuthMiddleware to protect routes
		protectedRoutes := superuserRoutes.Group("/")
		protectedRoutes.Use(middlewares.AuthTokenMiddleware(tokenManager))
		{
			// Protected routes
			protectedRoutes.GET("/dashboard", superuserHandler.DashboardSuperuserHandler)
			protectedRoutes.GET("/logout", superuserHandler.LogoutSuperuserHandler)
			protectedRoutes.GET("/test", superuserHandler.TestTemplate)

			// Profile routes
			protectedRoutes.GET("/profile", superuserHandler.ProfileViewHandler)
			protectedRoutes.POST("/profile", superuserHandler.ProfileUpdateHandler)

			// Password reset routes
			protectedRoutes.GET("/password-reset-request", superuserHandler.PasswordResetRequestHandler)
			protectedRoutes.POST("/password-reset/:token", superuserHandler.PasswordResetHandler)

			// 2FA routes
			protectedRoutes.GET("/enable-2fa", superuserHandler.Enable2FAHandler)
			protectedRoutes.POST("/verify-2fa", superuserHandler.Verify2FAHandler)

			// File upload and download
			protectedRoutes.POST("/upload", superuserHandler.FileUploadHandler)
			protectedRoutes.GET("/download/:filename", superuserHandler.FileDownloadHandler)
		}
	}
}
