package implement

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"gorm.io/gorm"
)

type serviceRepoImpl struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) repository.ServiceRepository {
	return &serviceRepoImpl{db}
}

func (r *serviceRepoImpl) CreateServiceType(ctx context.Context, serviceType *model.ServiceType) error {
	return r.db.WithContext(ctx).Create(serviceType).Error
}
