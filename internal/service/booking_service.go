package service

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

type BookingService interface {
	GetBookings(ctx context.Context, query types.BookingPaginationQuery) ([]*model.Booking, *types.MetaResponse, error)

	GetBookingByID(ctx context.Context, id int64) (*model.Booking, error)
}