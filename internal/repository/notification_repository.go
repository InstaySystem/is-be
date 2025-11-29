package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
	"gorm.io/gorm"
)

type Notification interface {
	CreateNotificationTx(tx *gorm.DB, notification *model.Notification) error

	FindAllUnreadNotificationsByContentIDAndType(ctx context.Context, staffID, contentID int64, contentType string) ([]*model.Notification, error)

	CreateNotificationStaffsTx(tx *gorm.DB, notificationStaffs []*model.NotificationStaff) error

	CreateNotificationStaffs(ctx context.Context, notificationStaffs []*model.NotificationStaff) error

	FindAllUnreadNotificationsByDepartmentID(ctx context.Context, staffID, departmentID int64) ([]*model.Notification, error)

	FindAllUnreadNotificationIDsByDepartmentIDTx(tx *gorm.DB, staffID, departmentID int64) ([]int64, error)

	FindAllUnreadNotificationIDsByOrderRoomIDTx(tx *gorm.DB, orderRoomID int64) ([]int64, error)

	FindAllNotificationsByOrderRoomIDTx(tx *gorm.DB, orderRoomID int64) ([]*model.Notification, error)

	CountUnreadNotificationsByDepartmentID(ctx context.Context, userID, departmentID int64) (int64, error)

	CountUnreadNotificationsByOrderRoomID(ctx context.Context, orderRoomID int64) (int64, error)

	UpdateReadNotificationsByOrderRoomTx(tx *gorm.DB, orderRoomID int64, updateData map[string]any) error

	UpdateNotificationsByOrderRoomIDAndType(ctx context.Context, orderRoomID int64, contentType string, updateData map[string]any) error

	FindAllNotificationsByDepartmentIDWithStaffsReadPaginatedTx(tx *gorm.DB, query types.NotificationPaginationQuery, staffID, departmentID int64) ([]*model.Notification, int64, error)
}
