package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
)

type ServiceRepository interface {
	CreateServiceType(ctx context.Context, service_type *model.ServiceType) error
}