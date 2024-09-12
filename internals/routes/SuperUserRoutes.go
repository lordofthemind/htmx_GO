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
		superuserRoutes.GET("/", superuserHandler.Index)
		superuserRoutes.GET("/register", superuserHandler.Register)
		superuserRoutes.GET("/login", superuserHandler.Login)
		superuserRoutes.POST("/register", superuserHandler.RegisterSuperuser)
		superuserRoutes.POST("/login", superuserHandler.LoginSuperuser)

		// Apply JWTAuthMiddleware to protect routes
		protectedRoutes := superuserRoutes.Group("/")
		protectedRoutes.Use(middlewares.JWTAuthMiddleware(tokenManager))
		{
			protectedRoutes.GET("/test", superuserHandler.TestTemplate)
			protectedRoutes.GET("/dashboard", superuserHandler.Dashboard)
			protectedRoutes.GET("/logout", superuserHandler.LogoutSuperuser)
		}
	}
}
