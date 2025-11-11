package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
)

type DepartmentRepository interface {
	Create(ctx context.Context, department *model.Department) error

	FindAllWithCreatedByAndUpdatedBy(ctx context.Context) ([]*model.Department, error)
}