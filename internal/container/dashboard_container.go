package container

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/repository"
	svcImpl "github.com/InstaySystem/is-be/internal/service/implement"
	"go.uber.org/zap"
)

type DashboardContainer struct {
	Hdl *handler.DashboardHandler
}

func NewDashboardContainer(
	userRepo repository.UserRepository,
	roomRepo repository.RoomRepository,
	serviceRepo repository.ServiceRepository,
	bookingRepo repository.BookingRepository,
	orderRepo repository.OrderRepository,
	requestRepo repository.RequestRepository,
	reviewRepo repository.ReviewRepository,
	logger      *zap.Logger,
) *DashboardContainer {
	svc := svcImpl.NewDashboardService(userRepo, roomRepo, serviceRepo, bookingRepo, orderRepo, requestRepo, reviewRepo, logger)
	hdl := handler.NewDashboardHandler(svc)

	return &DashboardContainer{hdl}
}
