package service

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

type NotificationService interface {
	GetNotificationsForAdmin(ctx context.Context, query types.NotificationPaginationQuery, userID, departmentID int64) ([]*model.Notification, *types.MetaResponse, error)

	CountUnreadNotificationsForAdmin(ctx context.Context, userID, departmentID int64) (int64, error)

	GetNotificationsForGuest(ctx context.Context, orderRoomID int64) ([]*model.Notification, error)

	CountUnreadNotificationsForGuest(ctx context.Context, orderRoomID int64) (int64, error)
}