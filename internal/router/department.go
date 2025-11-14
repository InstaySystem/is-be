package router

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func DepartmentRouter(rg *gin.RouterGroup, hdl *handler.DepartmentHandler, authMid *middleware.AuthMiddleware) {
	admin := rg.Group("/admin", authMid.IsAuthentication(), authMid.HasAnyRole([]string{"admin"}))
	{
		admin.POST("/departments", hdl.CreateDepartment)

		admin.GET("/departments", hdl.GetDepartments)

		admin.PATCH("/departments/:id", hdl.UpdateDepartment)

		admin.DELETE("/departments/:id", hdl.DeleteDepartment)
	}

	rg.GET("/departments", hdl.GetSimpleDepartments)
}