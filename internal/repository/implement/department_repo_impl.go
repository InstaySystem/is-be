package implement

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"gorm.io/gorm"
)

type departmentRepoImpl struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) repository.DepartmentRepository {
	return &departmentRepoImpl{db}
}

func (r *departmentRepoImpl) Create(ctx context.Context, department *model.Department) error {
	return r.db.WithContext(ctx).Create(department).Error
}

func (r *departmentRepoImpl) FindAllWithCreatedByAndUpdatedBy(ctx context.Context) ([]*model.Department, error) {
	var departments []*model.Department
	if err := r.db.WithContext(ctx).Preload("CreatedBy").Preload("UpdatedBy").Order("name ASC").Find(&departments).Error; err != nil {
		return nil, err
	}

	return departments, nil
}