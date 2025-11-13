package implement

import (
	"context"
	"errors"
	"time"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"go.uber.org/zap"
)

type serviceSvcImpl struct {
	serviceRepo repository.ServiceRepository
	sfGen       snowflake.Generator
	logger      *zap.Logger
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
		s.logger.Error("generate service type ID failed", zap.Error(err))
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

func (s *serviceSvcImpl) GetServiceTypesForAdmin(ctx context.Context) ([]*model.ServiceType, error) {
	serviceTypes, err := s.serviceRepo.FindAllServiceTypesWithDetails(ctx)
	if err != nil {
		s.logger.Error("get service types for admin failed", zap.Error(err))
		return nil, err
	}

	return serviceTypes, nil
}

func (s *serviceSvcImpl) UpdateServiceType(ctx context.Context, serviceTypeID, userID int64, req types.UpdateServiceTypeRequest) error {
	serviceType, err := s.serviceRepo.FindServiceTypeByID(ctx, serviceTypeID)
	if err != nil {
		s.logger.Error("find service type by id failed", zap.Int64("id", serviceTypeID), zap.Error(err))
		return err
	}
	if serviceType == nil {
		return common.ErrServiceTypeNotFound
	}

	updateData := map[string]any{}

	if req.Name != nil && serviceType.Name != *req.Name {
		updateData["name"] = req.Name
		updateData["slug"] = common.GenerateSlug(*req.Name)
	}
	if req.DepartmentID != nil && serviceType.DepartmentID != *req.DepartmentID {
		updateData["department_id"] = req.DepartmentID
	}

	if len(updateData) > 0 {
		updateData["updated_by_id"] = userID
		if err := s.serviceRepo.UpdateServiceType(ctx, serviceTypeID, updateData); err != nil {
			if ok, _ := common.IsUniqueViolation(err); ok {
				return common.ErrServiceTypeAlreadyExists
			}
			if common.IsForeignKeyViolation(err) {
				return common.ErrDepartmentNotFound
			}
			s.logger.Error("update service type failed", zap.Int64("id", serviceTypeID), zap.Error(err))
			return err
		}
	}

	return nil
}

func (s *serviceSvcImpl) DeleteServiceType(ctx context.Context, serviceTypeID int64) error {
	if err := s.serviceRepo.DeleteServiceType(ctx, serviceTypeID); err != nil {
		if errors.Is(err, common.ErrServiceTypeNotFound) {
			return err
		}
		if common.IsForeignKeyViolation(err) {
			return common.ErrProtectedRecord
		}
		s.logger.Error("delete service type failed", zap.Int64("id", serviceTypeID), zap.Error(err))
		return err
	}

	return nil
}

func (s *serviceSvcImpl) CreateService(ctx context.Context, userID int64, req types.CreateServiceRequest) (int64, error) {
	serviceID, err := s.sfGen.NextID()
	if err != nil {
		s.logger.Error("generate service ID failed", zap.Error(err))
		return 0, err
	}

	service := &model.Service{
		ID:            serviceID,
		Name:          req.Name,
		Slug:          common.GenerateSlug(req.Name),
		Price:         req.Price,
		CreatedByID:   userID,
		UpdatedByID:   userID,
		ServiceTypeID: req.ServiceTypeID,
	}

	serviceImages := make([]*model.ServiceImage, 0, len(req.Images))
	for _, reqImg := range req.Images {
		imageID, err := s.sfGen.NextID()
		if err != nil {
			s.logger.Error("generate service image ID failed", zap.Error(err))
			return 0, err
		}
		serviceImage := &model.ServiceImage{
			ID:          imageID,
			ServiceID:   serviceID,
			Key:         reqImg.Key,
			IsThumbnail: reqImg.IsThumbnail,
			SortOrder:   reqImg.SortOrder,
			UploadedAt:  time.Now(),
		}

		serviceImages = append(serviceImages, serviceImage)
	}

	service.ServiceImages = serviceImages

	if err = s.serviceRepo.CreateService(ctx, service); err != nil {
		if ok, _ := common.IsUniqueViolation(err); ok {
			return 0, common.ErrServiceAlreadyExists
		}
		if common.IsForeignKeyViolation(err) {
			return 0, common.ErrServiceTypeNotFound
		}
		s.logger.Error("create service failed", zap.Error(err))
		return 0, err
	}

	return serviceID, nil
}

func (s *serviceSvcImpl) GetServicesForAdmin(ctx context.Context, query types.ServicePaginationQuery) ([]*model.Service, *types.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	services, total, err := s.serviceRepo.FindAllServicesWithServiceTypeAndThumbnailPaginated(ctx, query)
	if err != nil {
		s.logger.Error("find all services paginated failed", zap.Error(err))
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

	return services, meta, nil
}
