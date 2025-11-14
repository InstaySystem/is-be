package router

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func UserRouter(rg *gin.RouterGroup, hdl *handler.UserHandler, authMid *middleware.AuthMiddleware) {
	user := rg.Group("/admin/users", authMid.IsAuthentication(), authMid.HasAnyRole([]string{"admin"}))
	{
		user.POST("", hdl.CreateUser)

		user.GET("/:id", hdl.GetUserByID)

		user.GET("", hdl.GetUsers)

		user.GET("/roles", hdl.GetAllRoles)

		user.PATCH("/:id", hdl.UpdateUser)

		user.PUT("/:id/password", hdl.UpdateUserPassword)

		user.DELETE("/:id", hdl.DeleteUser)
	}
}
