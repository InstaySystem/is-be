package implement

import (
	"context"
	"time"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type notificationSvcImpl struct {
	db               *gorm.DB
	notificationRepo repository.Notification
	logger           *zap.Logger
	sfGen            snowflake.Generator
}

func NewNotificationService(
	db *gorm.DB,
	notificationRepo repository.Notification,
	logger *zap.Logger,
	sfGen snowflake.Generator,
) service.NotificationService {
	return &notificationSvcImpl{
		db,
		notificationRepo,
		logger,
		sfGen,
	}
}

func (s *notificationSvcImpl) GetNotificationsForAdmin(ctx context.Context, query types.NotificationPaginationQuery, userID, departmentID int64) ([]*model.Notification, *types.MetaResponse, error) {
	var notifications []*model.Notification
	var total int64

	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		unreadNotificationIDs, err := s.notificationRepo.FindAllUnreadNotificationIDsByDepartmentIDTx(tx, userID, departmentID)
		if err != nil {
			s.logger.Error("find all unread notification ids by department id failed", zap.Error(err))
			return err
		}

		if len(unreadNotificationIDs) > 0 {
			notificationStaffs := make([]*model.NotificationStaff, 0, len(unreadNotificationIDs))
			for _, notificationID := range unreadNotificationIDs {
				id, err := s.sfGen.NextID()
				if err != nil {
					s.logger.Error("generate notification staff id failed", zap.Error(err))
					return err
				}

				notificationStaffs = append(notificationStaffs, &model.NotificationStaff{
					ID:             id,
					NotificationID: notificationID,
					StaffID:        userID,
				})
			}

			if err = s.notificationRepo.CreateNotificationStaffsTx(tx, notificationStaffs); err != nil {
				s.logger.Error("create notification staffs failed", zap.Error(err))
				return err
			}
		}

		notifications, total, err = s.notificationRepo.FindAllNotificationsByDepartmentIDWithStaffsReadPaginatedTx(tx, query, userID, departmentID)
		if err != nil {
			s.logger.Error("find all notifications by department id paginated failed", zap.Error(err))
			return err
		}

		return nil
	}); err != nil {
		return nil, nil, err
	}

	totalPages := uint32(total) / query.Limit
	if uint32(total)%query.Limit != 0 {
		totalPages++
	}

	meta := &types.MetaResponse{
		Total:      uint64(total),
		Page:       query.Page,
		Limit:      query.Limit,
		TotalPages: uint16(totalPages),
		HasPrev:    query.Page > 1,
		HasNext:    query.Page < totalPages,
	}

	return notifications, meta, nil
}

func (s *notificationSvcImpl) CountUnreadNotificationsForAdmin(ctx context.Context, userID, departmentID int64) (int64, error) {
	count, err := s.notificationRepo.CountUnreadNotificationsByDepartmentID(ctx, userID, departmentID)
	if err != nil {
		s.logger.Error("count unread notifications by department id failed", zap.Error(err))
		return 0, err
	}

	return count, nil
}

func (s *notificationSvcImpl) GetNotificationsForGuest(ctx context.Context, orderRoomID int64) ([]*model.Notification, error) {
	var notifications []*model.Notification
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		unreadNotificationIDs, err := s.notificationRepo.FindAllUnreadNotificationIDsByOrderRoomIDTx(tx, orderRoomID)
		if err != nil {
			s.logger.Error("find all unread notification ids by order room id failed", zap.Error(err))
			return err
		}

		if len(unreadNotificationIDs) > 0 {
			notificationIDs := []int64{}
			for _, notificationID := range unreadNotificationIDs {
				notificationIDs = append(notificationIDs, notificationID)
			}

			updateData := map[string]any{
				"is_read": true,
				"read_at": time.Now(),
			}
			if err = s.notificationRepo.UpdateReadNotificationsByOrderRoomTx(tx, orderRoomID, updateData); err != nil {
				s.logger.Error("update read notifications failed", zap.Error(err))
				return err
			}
		}

		notifications, err = s.notificationRepo.FindAllNotificationsByOrderRoomIDTx(tx, orderRoomID)
		if err != nil {
			s.logger.Error("find all notifications by order room id failed", zap.Error(err))
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return notifications, nil
}

func (s *notificationSvcImpl) CountUnreadNotificationsForGuest(ctx context.Context, orderRoomID int64) (int64, error) {
	count, err := s.notificationRepo.CountUnreadNotificationsByOrderRoomID(ctx, orderRoomID)
	if err != nil {
		s.logger.Error("count unread notifications by order room id failed", zap.Error(err))
		return 0, err
	}

	return count, nil
}
