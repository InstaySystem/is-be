package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/config"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authSvc service.AuthService
	cfg     *config.Config
}

func NewAuthHandler(authSvc service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authSvc,
		cfg,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req types.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	user, accessToken, refreshToken, err := h.authSvc.Login(ctx, req)
	if err != nil {
		switch err {
		case common.ErrLoginFailed:
			common.ToAPIResponse(c, http.StatusBadRequest, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	c.SetCookie(h.cfg.JWT.AccessName, accessToken, int(h.cfg.JWT.AccessExpiresIn.Seconds()), "/", "", false, true)
	c.SetCookie(h.cfg.JWT.RefreshName, refreshToken, int(h.cfg.JWT.RefreshExpiresIn.Seconds()), fmt.Sprintf("%s/auth/refresh", h.cfg.Server.APIPrefix), "", false, true)

	common.ToAPIResponse(c, http.StatusOK, "Login successfully", gin.H{
		"user": common.ToUserResponse(user),
	})
}
