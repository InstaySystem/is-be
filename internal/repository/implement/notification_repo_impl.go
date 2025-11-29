package implement

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/types"
	"gorm.io/gorm"
)

type notificationRepoImpl struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) repository.Notification {
	return &notificationRepoImpl{db}
}

func (r *notificationRepoImpl) CreateNotificationTx(tx *gorm.DB, notification *model.Notification) error {
	return tx.Create(notification).Error
}

func (r *notificationRepoImpl) CreateNotificationStaffsTx(tx *gorm.DB, notificationStaffs []*model.NotificationStaff) error {
	return tx.Create(notificationStaffs).Error
}

func (r *notificationRepoImpl) CreateNotificationStaffs(ctx context.Context, notificationStaffs []*model.NotificationStaff) error {
	return r.db.WithContext(ctx).Create(notificationStaffs).Error
}

func (r *notificationRepoImpl) UpdateReadNotificationsByOrderRoomTx(tx *gorm.DB, orderRoomID int64, updateData map[string]any) error {
	return tx.Model(&model.Notification{}).Where("order_room_id = ? AND receiver = ? AND is_read = false", orderRoomID, "guest").Updates(updateData).Error
}

func (r *notificationRepoImpl) FindAllUnreadNotificationsByContentIDAndType(ctx context.Context, staffID, contentID int64, contentType string) ([]*model.Notification, error) {
	var notifications []*model.Notification
	if err := r.db.WithContext(ctx).Where("content_id = ? AND type = ? AND receiver = ?", contentID, contentType, "staff").Where("id NOT IN (?)",
		r.db.Model(&model.NotificationStaff{}).
			Select("notification_id").
			Where("staff_id = ?", staffID),
	).Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *notificationRepoImpl) FindAllNotificationsByOrderRoomIDTx(tx *gorm.DB, orderRoomID int64) ([]*model.Notification, error) {
	var notifications []*model.Notification
	if err := tx.Model(&model.Notification{}).Where("order_room_id = ? AND receiver = ?", orderRoomID, "guest").Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *notificationRepoImpl) FindAllUnreadNotificationIDsByOrderRoomIDTx(tx *gorm.DB, orderRoomID int64) ([]int64, error) {
	var ids []int64
	if err := tx.Model(&model.Notification{}).Where("order_room_id = ? AND receiver = ? AND is_read = false", orderRoomID, "guest").Pluck("id", &ids).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *notificationRepoImpl) FindAllUnreadNotificationIDsByDepartmentIDTx(tx *gorm.DB, staffID, departmentID int64) ([]int64, error) {
	var ids []int64
	if err := tx.Where("department_id = ? AND receiver = ?", departmentID, "staff").Where("id NOT IN (?)",
		tx.Model(&model.NotificationStaff{}).
			Select("notification_id").
			Where("staff_id = ?", staffID),
	).Model(&model.Notification{}).Pluck("id", &ids).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *notificationRepoImpl) FindAllUnreadNotificationsByDepartmentID(ctx context.Context, staffID, departmentID int64) ([]*model.Notification, error) {
	var notifications []*model.Notification
	if err := r.db.WithContext(ctx).Where("department_id = ? AND receiver = ?", departmentID, "staff").Where("id NOT IN (?)",
		r.db.Model(&model.NotificationStaff{}).
			Select("notification_id").
			Where("staff_id = ?", staffID),
	).Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *notificationRepoImpl) CountUnreadNotificationsByDepartmentID(ctx context.Context, userID, departmentID int64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Notification{}).Where("department_id = ? AND receiver = ?", departmentID, "staff").
		Where("id NOT IN (?)",
			r.db.Model(&model.NotificationStaff{}).
				Select("notification_id").
				Where("staff_id = ?", userID),
		).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *notificationRepoImpl) CountUnreadNotificationsByOrderRoomID(ctx context.Context, orderRoomID int64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Notification{}).Where("order_room_id = ? AND receiver = ? AND is_read = false", orderRoomID, "guest").Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *notificationRepoImpl) UpdateNotificationsByOrderRoomIDAndType(ctx context.Context, orderRoomID int64, contentType string, updateData map[string]any) error {
	return r.db.WithContext(ctx).Model(&model.Notification{}).Where("order_room_id = ? AND type = ? AND receiver = ?", orderRoomID, contentType, "guest").Updates(updateData).Error
}

func (r *notificationRepoImpl) FindAllNotificationsByDepartmentIDWithStaffsReadPaginatedTx(tx *gorm.DB, query types.NotificationPaginationQuery, staffID, departmentID int64) ([]*model.Notification, int64, error) {
	var notifications []*model.Notification
	var total int64

	db := tx.Where("department_id = ? AND receiver = ?", departmentID, "staff").Model(&model.Notification{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.Limit
	if err := db.Preload("StaffsRead", "staff_id = ?", staffID).Order("created_at DESC").Limit(int(query.Limit)).Offset(int(offset)).Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}
