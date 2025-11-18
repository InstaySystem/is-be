package service

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

type RoomService interface {
	CreateRoomType(ctx context.Context, userID int64, req types.CreateRoomTypeRequest) error

	GetRoomTypesForAdmin(ctx context.Context) ([]*model.RoomType, error)
}