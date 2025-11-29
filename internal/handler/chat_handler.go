package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatSvc service.ChatService
}

func NewChatHandler(chatSvc service.ChatService) *ChatHandler {
	return &ChatHandler{chatSvc}
}

func (h *ChatHandler) GetChatsForAdmin(c *gin.Context) {
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

	var query types.ChatPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	if user.Department == nil {
		common.ToAPIResponse(c, http.StatusOK, "Get chat list successfully", gin.H{
			"chats": []any{},
			"meta":  &types.MetaResponse{},
		})
		return
	}

	chats, meta, err := h.chatSvc.GetChatsForAdmin(ctx, query, user.ID, user.Department.ID)
	if err != nil {
		common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get chat list successfully", gin.H{
		"chats": common.ToSimpleChatsResponse(chats),
		"meta":  meta,
	})
}

func (h *ChatHandler) GetChatByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	chatIDStr := c.Param("id")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		common.ToAPIResponse(c, http.StatusBadRequest, common.ErrInvalidID.Error(), nil)
		return
	}

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

	chat, err := h.chatSvc.GetChatByID(ctx, chatID, user.ID, user.Department.ID)
	if err != nil {
		switch err {
		case common.ErrChatNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		case common.ErrForbidden:
			common.ToAPIResponse(c, http.StatusForbidden, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get chat information successfully", gin.H{
		"chat": common.ToSimpleChatWithMessagesResponse(chat),
	})
}

func (h *ChatHandler) GetChatsForGuest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	orderRoomID := c.GetInt64("order_room_id")
	if orderRoomID == 0 {
		common.ToAPIResponse(c, http.StatusForbidden, common.ErrForbidden.Error(), nil)
		return
	}

	chats, err := h.chatSvc.GetChatsForGuest(ctx, orderRoomID)
	if err != nil {
		common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get chat list successfully", gin.H{
		"chats": common.ToBasicChatsResponse(chats),
	})
}

func (h *ChatHandler) GetChatByCode(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	chatCode := c.Param("code")

	orderRoomID := c.GetInt64("order_room_id")
	if orderRoomID == 0 {
		common.ToAPIResponse(c, http.StatusForbidden, common.ErrForbidden.Error(), nil)
		return
	}

	chat, err := h.chatSvc.GetChatByCode(ctx, chatCode, orderRoomID)
	if err != nil {

	}
	common.ToAPIResponse(c, http.StatusForbidden, common.ErrForbidden.Error(), chat)
}
