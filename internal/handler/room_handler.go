package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	roomSvc service.RoomService
}

func NewRoomHandler(roomSvc service.RoomService) *RoomHandler {
	return &RoomHandler{roomSvc}
}

func (h *RoomHandler) CreateRoomType(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.ToAPIResponse(c, http.StatusUnauthorized, common.ErrUnAuth.Error(), nil)
		return
	}

	user, ok := userAny.(*types.UserData)
	if !ok {
		common.ToAPIResponse(c, http.StatusUnauthorized, common.ErrInvalidUser.Error(), nil)
		return
	}

	var req types.CreateRoomTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	if err := h.roomSvc.CreateRoomType(ctx, user.ID, req); err != nil {
		switch err {
		case common.ErrRoomTypeAlreadyExists:
			common.ToAPIResponse(c, http.StatusConflict, common.ErrUnAuth.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusCreated, "Room type created successfully", nil)
}

func (h *RoomHandler) GetRoomTypesForAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	roomTypes, err := h.roomSvc.GetRoomTypesForAdmin(ctx)
	if err != nil {
		common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get room types successfully", gin.H{
		"room_types": common.ToRoomTypesResponse(roomTypes),
	})
}
