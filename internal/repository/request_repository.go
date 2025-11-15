package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
)

type RequestRepository interface {
	CreateRequestType(ctx context.Context, requestType *model.RequestType) error

	FindAllRequestTypeWithDetails(ctx context.Context) ([]*model.RequestType, error)

	FindRequestTypeByID(ctx context.Context, requestTypeID int64) (*model.RequestType, error)

	UpdateRequestType(ctx context.Context, requestTypeID int64, updateData map[string]any) error

	DeleteRequestType(ctx context.Context, requestTypeID int64) error
}