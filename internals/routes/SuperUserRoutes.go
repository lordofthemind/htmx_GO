package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lordofthemind/htmx_GO/internals/handlers"
	"github.com/lordofthemind/htmx_GO/internals/middlewares"
	"github.com/lordofthemind/htmx_GO/internals/tokens"
)

func RegisterSuperuserRoutes(router *gin.Engine, superuserHandler *handlers.SuperuserHandler, tokenManager tokens.TokenManager) {
	superuserRoutes := router.Group("/superuser")
	{
		superuserRoutes.GET("/", superuserHandler.IndexRender)
		superuserRoutes.GET("/register", superuserHandler.RegisterRender)
		superuserRoutes.GET("/login", superuserHandler.LoginRender)
		superuserRoutes.POST("/register", superuserHandler.RegisterSuperuserHandler)
		superuserRoutes.POST("/login", superuserHandler.LoginSuperuserHandler)

		// Apply JWTAuthMiddleware to protect routes
		protectedRoutes := superuserRoutes.Group("/")
		protectedRoutes.Use(middlewares.JWTAuthMiddleware(tokenManager))
		{
			protectedRoutes.GET("/dashboard", superuserHandler.DashboardSuperuserHandler)
			protectedRoutes.GET("/logout", superuserHandler.LogoutSuperuserHandler)
			protectedRoutes.GET("/test", superuserHandler.TestTemplate)
		}
	}
}
