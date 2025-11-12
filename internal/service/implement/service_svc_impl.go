package implement

import (
	"context"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"go.uber.org/zap"
)

type serviceSvcImpl struct {
	serviceRepo    repository.ServiceRepository
	sfGen          snowflake.Generator
	logger         *zap.Logger
}

func NewServiceService(
	serviceRepo repository.ServiceRepository,
	sfGen snowflake.Generator,
	logger *zap.Logger,
) service.ServiceService {
	return &serviceSvcImpl{
		serviceRepo,
		sfGen,
		logger,
	}
}

func (s *serviceSvcImpl) CreateServiceType(ctx context.Context, userID int64, req types.CreateServiceTypeRequest) error {
	id, err := s.sfGen.NextID()
	if err != nil {
		s.logger.Error("generate ID failed", zap.Error(err))
		return err
	}

	serviceType := &model.ServiceType{
		ID:           id,
		Name:         req.Name,
		Slug:         common.GenerateSlug(req.Name),
		DepartmentID: req.DepartmentID,
		CreatedByID:  userID,
		UpdatedByID:  userID,
	}

	if err = s.serviceRepo.CreateServiceType(ctx, serviceType); err != nil {
		if ok, _ := common.IsUniqueViolation(err); ok {
			return common.ErrServiceTypeAlreadyExists
		}
		if common.IsForeignKeyViolation(err) {
			return common.ErrDepartmentNotFound
		}
		s.logger.Error("create service type failed", zap.Error(err))
		return err
	}

	return nil
}
