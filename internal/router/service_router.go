package router

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func ServiceRouter(rg *gin.RouterGroup, hdl *handler.ServiceHandler, authMid *middleware.AuthMiddleware) {
	admin := rg.Group("/admin", authMid.IsAuthentication(), authMid.HasAnyRole([]string{"admin"}))
	{
		admin.POST("/service-types", hdl.CreateServiceType)

		admin.PATCH("/service-types/:id", hdl.UpdateServiceType)

		admin.DELETE("/service-types/:id", hdl.DeleteServiceType)

		admin.POST("/services", hdl.CreateService)

		admin.PATCH("/services/:id", hdl.UpdateService)

		admin.DELETE("/services/:id", hdl.DeleteService)
	}

	admin = rg.Group("/admin", authMid.IsAuthentication(), authMid.HasDepartment("reception"))
	{
		admin.GET("/services", hdl.GetServicesForAdmin)

		admin.GET("/services/:id", hdl.GetServiceByID)

		admin.GET("/service-types", hdl.GetServiceTypesForAdmin)
	}

	rg.GET("/service-types", hdl.GetServiceTypesForGuest)

	rg.GET("/service-types/:slug/services", authMid.HasGuestToken(), hdl.GetServiceTypeBySlugWithServices)

	rg.GET("/services/:slug", authMid.HasGuestToken(), hdl.GetServiceBySlug)
}
