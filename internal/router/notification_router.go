package router

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func NotificationRouter(rg *gin.RouterGroup, hdl *handler.NotificationHandler, authMid *middleware.AuthMiddleware) {
	admin := rg.Group("/admin/notifications", authMid.IsAuthentication())
	{
		admin.GET("", hdl.GetNotificationsForAdmin)

		admin.GET("/unread-count", hdl.CountUnreadNotificationsForAdmin)
	}

	guest := rg.Group("/notifications", authMid.HasGuestToken())
	{
		guest.GET("", hdl.GetNotificationsForGuest)

		guest.GET("/unread-count", hdl.CountUnreadNotificationsForGuest)
	}
}