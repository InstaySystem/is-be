package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

type ServiceRepository interface {
	CreateServiceType(ctx context.Context, serviceType *model.ServiceType) error

	FindAllServiceTypesWithDetails(ctx context.Context) ([]*model.ServiceType, error)

	FindServiceTypeByID(ctx context.Context, serviceTypeID int64) (*model.ServiceType, error)

	UpdateServiceType(ctx context.Context, serviceTypeID int64, updateData map[string]any) error

	DeleteServiceType(ctx context.Context, serviceTypeID int64) error

	CreateService(ctx context.Context, service *model.Service) error

	FindAllServicesWithServiceTypeAndThumbnailPaginated(ctx context.Context, query types.ServicePaginationQuery) ([]*model.Service, int64, error)
}