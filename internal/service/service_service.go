package service

import (
	"context"

	"github.com/InstaySystem/is-be/internal/types"
)

type ServiceService interface {
	CreateServiceType(ctx context.Context, userID int64, req types.CreateServiceTypeRequest) error
}