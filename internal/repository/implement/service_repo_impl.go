package implement

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/types"
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

func (r *serviceRepoImpl) FindAllServiceTypesWithDetails(ctx context.Context) ([]*model.ServiceType, error) {
	var serviceTypes []*model.ServiceType
	if err := r.db.WithContext(ctx).Preload("Department").Preload("CreatedBy").Preload("UpdatedBy").Order("name ASC").Find(&serviceTypes).Error; err != nil {
		return nil, err
	}

	return serviceTypes, nil
}

func (r *serviceRepoImpl) FindServiceTypeByID(ctx context.Context, serviceTypeID int64) (*model.ServiceType, error) {
	var serviceType model.ServiceType
	if err := r.db.WithContext(ctx).Where("id = ?", serviceTypeID).First(&serviceType).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &serviceType, nil
}

func (r *serviceRepoImpl) UpdateServiceType(ctx context.Context, serviceTypeID int64, updateData map[string]any) error {
	result := r.db.WithContext(ctx).Model(&model.ServiceType{}).Where("id = ?", serviceTypeID).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.ErrServiceTypeNotFound
	}

	return nil
}

func (r *serviceRepoImpl) DeleteServiceType(ctx context.Context, serviceTypeID int64) error {
	result := r.db.WithContext(ctx).Where("id = ?", serviceTypeID).Delete(&model.ServiceType{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.ErrServiceTypeNotFound
	}

	return nil
}

func (r *serviceRepoImpl) CreateService(ctx context.Context, service *model.Service) error {
	return r.db.WithContext(ctx).Create(service).Error
}

func (r *serviceRepoImpl) FindAllServicesWithServiceTypeAndThumbnailPaginated(ctx context.Context, query types.ServicePaginationQuery) ([]*model.Service, int64, error) {
	var services []*model.Service
	var total int64

	db := r.db.WithContext(ctx).Preload("ServiceType").Preload("ServiceImages", "is_thumbnail = true").Model(&model.Service{})
	db = applyServiceFilters(db, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyServiceSorting(db, query)
	offset := (query.Page - 1) * query.Limit
	if err := db.Offset(int(offset)).Limit(int(query.Limit)).Find(&services).Error; err != nil {
		return nil, 0, err
	}

	return services, total, nil
}

func applyServiceFilters(db *gorm.DB, query types.ServicePaginationQuery) *gorm.DB {
	if query.Search != "" {
		searchTerm := "%" + strings.ToLower(query.Search) + "%"
		db = db.Where(
			"LOWER(name) LIKE @q OR LOWER(slug) LIKE @q",
			sql.Named("q", searchTerm),
		)
	}

	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}

	if query.ServiceTypeID != 0 {
		db = db.Where("service_type_id = ?", query.ServiceTypeID)
	}

	return db
}

func applyServiceSorting(db *gorm.DB, query types.ServicePaginationQuery) *gorm.DB {
	if query.Sort == "" {
		query.Sort = "created_at"
	}
	if query.Order == "" {
		query.Order = "desc"
	}

	allowedSorts := map[string]bool{
		"created_at": true,
		"name":       true,
		"price":      true,
	}

	if allowedSorts[query.Sort] {
		db = db.Order(query.Sort + " " + strings.ToUpper(query.Order))
	} else {
		db = db.Order("created_at DESC")
	}

	return db
}
