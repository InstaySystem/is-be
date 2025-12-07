package service

import (
	"context"

	"github.com/InstaySystem/is-be/internal/types"
)

type DashboardService interface {
	Overview(ctx context.Context) (*types.DashboardResponse, error)
}