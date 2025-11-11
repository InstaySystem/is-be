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

type departmentSvcImpl struct {
	departmentRepo repository.DepartmentRepository
	sfGen          snowflake.Generator
	logger         *zap.Logger
}

func NewDepartmentService(
	departmentRepo repository.DepartmentRepository,
	sfGen snowflake.Generator,
	logger *zap.Logger,
) service.DepartmentService {
	return &departmentSvcImpl{
		departmentRepo,
		sfGen,
		logger,
	}
}

func (s *departmentSvcImpl) CreateDepartment(ctx context.Context, userID int64, req types.CreateDepartmentRequest) (int64, error) {
	id, err := s.sfGen.NextID()
	if err != nil {
		s.logger.Error("generate ID failed", zap.Error(err))
		return 0, err
	}

	department := &model.Department{
		ID:          id,
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		CreatedByID: userID,
		UpdatedByID: userID,
	}

	if err = s.departmentRepo.Create(ctx, department); err != nil {
		ok, _ := common.IsUniqueViolation(err)
		if ok {
			return 0, common.ErrDepartmentAlreadyExists
		}
		s.logger.Error("create department failed", zap.Error(err))
		return 0, err
	}

	return id, nil
}

func (s *departmentSvcImpl) GetDepartments(ctx context.Context) ([]*model.Department, error) {
	departments, err := s.departmentRepo.FindAllWithCreatedByAndUpdatedBy(ctx)
	if err != nil {
		s.logger.Error("get departments failed", zap.Error(err))
		return nil, err
	}

	return departments, nil
}
