package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

type BookingRepository interface {
	Create(ctx context.Context, booking *model.Booking) error

	FindAllPaginated(ctx context.Context, query types.BookingPaginationQuery) ([]*model.Booking, int64, error)

	FindByID(ctx context.Context, id int64) (*model.Booking, error)
}