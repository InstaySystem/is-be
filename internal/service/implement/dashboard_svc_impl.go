package implement

import (
	"context"
	"math"

	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type dashboardSvcImpl struct {
	userRepo    repository.UserRepository
	roomRepo    repository.RoomRepository
	serviceRepo repository.ServiceRepository
	bookingRepo repository.BookingRepository
	orderRepo   repository.OrderRepository
	requestRepo repository.RequestRepository
	logger      *zap.Logger
}

func NewDashboardService(
	userRepo repository.UserRepository,
	roomRepo repository.RoomRepository,
	serviceRepo repository.ServiceRepository,
	bookingRepo repository.BookingRepository,
	orderRepo repository.OrderRepository,
	requestRepo repository.RequestRepository,
	logger *zap.Logger,
) service.DashboardService {
	return &dashboardSvcImpl{
		userRepo,
		roomRepo,
		serviceRepo,
		bookingRepo,
		orderRepo,
		requestRepo,
		logger,
	}
}

func (s *dashboardSvcImpl) Overview(ctx context.Context) (*types.DashboardResponse, error) {
	var res types.DashboardResponse
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		count, err := s.userRepo.Count(ctx)
		if err != nil {
			return err
		}
		res.TotalStaff = count
		return nil
	})

	g.Go(func() error {
		count, err := s.roomRepo.CountRoom(ctx)
		if err != nil {
			return err
		}
		res.TotalRooms = count
		return nil
	})

	g.Go(func() error {
		count, err := s.serviceRepo.CountService(ctx)
		if err != nil {
			return err
		}
		res.TotalServices = count
		return nil
	})

	g.Go(func() error {
		count, err := s.bookingRepo.CountBooking(ctx)
		if err != nil {
			return err
		}
		res.TotalBookings = count
		return nil
	})

	g.Go(func() error {
		sum, err := s.bookingRepo.SumBookingTotalSellPrice(ctx)
		if err != nil {
			return err
		}
		res.BookingRevenue = sum
		return nil
	})

	g.Go(func() error {
		data, err := s.orderRepo.OrderServiceStatusDistribution(ctx)
		if err != nil {
			return err
		}

		calculatePercentage(data)
		res.OrderServiceStats = data
		return nil
	})

	g.Go(func() error {
		data, err := s.requestRepo.RequestStatusDistribution(ctx)
		if err != nil {
			return err
		}

		calculatePercentage(data)
		res.RequestStats = data
		return nil
	})

	if err := g.Wait(); err != nil {
		s.logger.Error("get dashboard failed", zap.Error(err))
		return nil, err
	}

	return &res, nil
}

func calculatePercentage(data []*types.StatusChartResponse) {
	var total int64
	for _, item := range data {
		total += item.Count
	}

	if total == 0 {
		return
	}

	for _, item := range data {
		item.Percentage = math.Round((float64(item.Count)/float64(total))*100*100) / 100
	}
}
