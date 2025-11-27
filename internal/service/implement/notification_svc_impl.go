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
)

type notificationSvcImpl struct {
	notificationRepo repository.Notification
	logger           *zap.Logger
	sfGen            snowflake.Generator
}

func NewNotificationService(
	notificationRepo repository.Notification,
	logger *zap.Logger,
	sfGen snowflake.Generator,
) service.NotificationService {
	return &notificationSvcImpl{
		notificationRepo,
		logger,
		sfGen,
	}
}

func (s *notificationSvcImpl) GetNotificationsForAdmin(ctx context.Context, query types.NotificationPaginationQuery, userID, departmentID int64) ([]*model.Notification, *types.MetaResponse, error) {
	unreadNotifications, err := s.notificationRepo.FindAllUnreadNotificationsByDepartmentID(ctx, userID, departmentID)
	if err != nil {
		s.logger.Error("find unread notifications by department ID failed", zap.Error(err))
		return nil, nil, err
	}

	if len(unreadNotifications) > 0 {
		notificationStaffs := make([]*model.NotificationStaff, 0, len(unreadNotifications))
		for _, notification := range unreadNotifications {
			id, err := s.sfGen.NextID()
			if err != nil {
				s.logger.Error("generate notification staff ID failed", zap.Error(err))
				return nil, nil, err
			}

			notificationStaffs = append(notificationStaffs, &model.NotificationStaff{
				ID:             id,
				NotificationID: notification.ID,
				StaffID:        userID,
			})
		}

		if err = s.notificationRepo.CreateNotificationStaffs(ctx, notificationStaffs); err != nil {
			s.logger.Error("create notification staffs failed", zap.Error(err))
			return nil, nil, err
		}
	}

	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	notifications, total, err := s.notificationRepo.FindAllNotificationsByDepartmentIDWithStaffsReadPaginated(ctx, query, userID, departmentID)
	if err != nil {
		s.logger.Error("find all notifications by department ID paginated failed", zap.Error(err))
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
		s.logger.Error("count unread notifications by department ID failed", zap.Error(err))
		return 0, err
	}

	return count, nil
}

func (s *notificationSvcImpl) GetNotificationsForGuest(ctx context.Context, orderRoomID int64) ([]*model.Notification, error) {
	unreadNotifications, err := s.notificationRepo.FindAllUnreadNotificationsByOrderRoomID(ctx, orderRoomID)
	if err != nil {
		s.logger.Error("find unread notifications by order room ID failed", zap.Error(err))
		return nil, err
	}

	if len(unreadNotifications) > 0 {
		notificationIDs := []int64{}
		for _, notification := range unreadNotifications {
			notificationIDs = append(notificationIDs, notification.ID)
		}

		updateData := map[string]any{
			"is_read": true,
			"read_at": time.Now(),
		}
		if err = s.notificationRepo.UpdateNotifications(ctx, notificationIDs, updateData); err != nil {
			s.logger.Error("update notifications failed", zap.Error(err))
			return nil, err
		}
	}

	notifications, err := s.notificationRepo.FindAllNotificationsByOrderRoomID(ctx, orderRoomID)
	if err != nil {
		s.logger.Error("find all notifications by order room ID failed", zap.Error(err))
		return nil, err
	}

	return notifications, nil
}

func (s *notificationSvcImpl) CountUnreadNotificationsForGuest(ctx context.Context, orderRoomID int64) (int64, error) {
	count, err := s.notificationRepo.CountUnreadNotificationsByOrderRoomID(ctx, orderRoomID)
	if err != nil {
		s.logger.Error("count unread notifications by order room ID failed", zap.Error(err))
		return 0, err
	}

	return count, nil
}
