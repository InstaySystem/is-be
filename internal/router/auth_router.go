package router

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/gin-gonic/gin"
)

func AuthRouter(rg *gin.RouterGroup, hdl *handler.AuthHandler) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", hdl.Login)
	}
} 