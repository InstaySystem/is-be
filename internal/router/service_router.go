package router

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func ServiceRouter(rg *gin.RouterGroup, hdl *handler.ServiceHandler, authMid *middleware.AuthMiddleware) {
	service_type := rg.Group("/service-types", authMid.IsAuthentication(), authMid.HasAnyRole([]string{"admin"}))
	{
		service_type.POST("", hdl.CreateServiceType)
	}
}