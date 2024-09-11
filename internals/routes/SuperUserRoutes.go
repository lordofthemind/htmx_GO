package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lordofthemind/htmx_GO/internals/handlers"
)

func RegisterSuperuserRoutes(router *gin.Engine, superuserHandler *handlers.SuperuserHandler) {
	superuserRoutes := router.Group("/superuser")
	{
		superuserRoutes.POST("/register", superuserHandler.RegisterSuperuser)
		superuserRoutes.POST("/login", superuserHandler.LoginSuperuser)
	}
}
