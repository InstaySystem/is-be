package implement

import (
	"context"
	"errors"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"go.uber.org/zap"
)

type requestSvcImpl struct {
	requestRepo repository.RequestRepository
	sfGen       snowflake.Generator
	logger      *zap.Logger
}

func NewRequestService(
	requestRepo repository.RequestRepository,
	sfGen snowflake.Generator,
	logger *zap.Logger,
) service.RequestService {
	return &requestSvcImpl{
		requestRepo,
		sfGen,
		logger,
	}
}

func (s *requestSvcImpl) CreateRequestType(ctx context.Context, userID int64, req types.CreateRequestTypeRequest) error {
	id, err := s.sfGen.NextID()
	if err != nil {
		s.logger.Error("generate request type ID failed", zap.Error(err))
		return err
	}

	requestType := &model.RequestType{
		ID:           id,
		Name:         req.Name,
		Slug:         common.GenerateSlug(req.Name),
		DepartmentID: req.DepartmentID,
		CreatedByID:  userID,
		UpdatedByID:  userID,
	}

	if err = s.requestRepo.CreateRequestType(ctx, requestType); err != nil {
		if ok, _ := common.IsUniqueViolation(err); ok {
			return common.ErrRequestTypeAlreadyExists
		}
		if common.IsForeignKeyViolation(err) {
			return common.ErrDepartmentNotFound
		}
		s.logger.Error("create request type failed", zap.Error(err))
		return err
	}

	return nil
}

func (s *requestSvcImpl) GetRequestTypesForAdmin(ctx context.Context) ([]*model.RequestType, error) {
	requestTypes, err := s.requestRepo.FindAllRequestTypeWithDetails(ctx)
	if err != nil {
		s.logger.Error("get request types for admin failed", zap.Error(err))
		return nil, err
	}

	return requestTypes, nil
}

func (s *requestSvcImpl) UpdateRequestType(ctx context.Context, requestTypeID, userID int64, req types.UpdateRequestTypeRequest) error {
	requestType, err := s.requestRepo.FindRequestTypeByID(ctx, requestTypeID)
	if err != nil {
		s.logger.Error("find request type by id failed", zap.Int64("id", requestTypeID), zap.Error(err))
		return err
	}
	if requestType == nil {
		return common.ErrRequestTypeNotFound
	}

	updateData := map[string]any{}

	if req.Name != nil && *req.Name != requestType.Name {
		updateData["name"] = *req.Name
		updateData["slug"] = common.GenerateSlug(*req.Name)
	}
	if req.DepartmentID != nil && *req.DepartmentID != requestType.DepartmentID {
		updateData["department_id"] = *req.DepartmentID
	}

	if len(updateData) > 0 {
		updateData["updated_by_id"] = userID
		if err := s.requestRepo.UpdateRequestType(ctx, requestTypeID, updateData); err != nil {
			if ok, _ := common.IsUniqueViolation(err); ok {
				return common.ErrRequestTypeAlreadyExists
			}
			if common.IsForeignKeyViolation(err) {
				return common.ErrDepartmentNotFound
			}
			s.logger.Error("update request type failed", zap.Int64("id", requestTypeID), zap.Error(err))
			return err
		}
	}

	return nil
}

func (s *requestSvcImpl) DeleteRequestType(ctx context.Context, requestTypeID int64) error {
	if err := s.requestRepo.DeleteRequestType(ctx, requestTypeID); err != nil {
		if errors.Is(err, common.ErrRequestTypeNotFound) {
			return err
		}
		if common.IsForeignKeyViolation(err) {
			return common.ErrProtectedRecord
		}
		s.logger.Error("delete request type failed", zap.Int64("id", requestTypeID), zap.Error(err))
		return err
	}

	return nil
}
