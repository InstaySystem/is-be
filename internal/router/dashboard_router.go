package router

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func DashboardRouter(rg *gin.RouterGroup, hdl *handler.DashboardHandler, authMid *middleware.AuthMiddleware) {
	admin := rg.Group("/admin/dashboard", authMid.IsAuthentication(), authMid.HasAnyRole([]string{"admin"}))
	{
		admin.GET("", hdl.Overview)
	}
}
