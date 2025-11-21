package implement

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/provider/cache"
	"github.com/InstaySystem/is-be/internal/provider/jwt"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"go.uber.org/zap"
)

type orderSvcImpl struct {
	orderRepo     repository.OrderRepository
	bookingRepo   repository.BookingRepository
	sfGen         snowflake.Generator
	logger        *zap.Logger
	cacheProvider cache.CacheProvider
	jwtProvider   jwt.JWTProvider
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	bookingRepo repository.BookingRepository,
	sfGen snowflake.Generator,
	logger *zap.Logger,
	cacheProvider cache.CacheProvider,
	jwtProvider jwt.JWTProvider,
) service.OrderService {
	return &orderSvcImpl{
		orderRepo,
		bookingRepo,
		sfGen,
		logger,
		cacheProvider,
		jwtProvider,
	}
}

func (s *orderSvcImpl) CreateOrderRoom(ctx context.Context, userID int64, req types.CreateOrderRoomRequest) (int64, string, error) {
	booking, err := s.bookingRepo.FindByID(ctx, req.BookingID)
	if err != nil {
		s.logger.Error("find booking by ID failed", zap.Int64("id", req.BookingID), zap.Error(err))
		return 0, "", err
	}

	if booking.CheckOut.Before(time.Now()) {
		return 0, "", common.ErrBookingExpired
	}

	id, err := s.sfGen.NextID()
	if err != nil {
		s.logger.Error("generate order room ID failed", zap.Error(err))
		return 0, "", err
	}

	orderRoom := &model.OrderRoom{
		ID:          id,
		CreatedByID: userID,
		UpdatedByID: userID,
		BookingID:   req.BookingID,
		RoomID:      req.RoomID,
	}

	if err = s.orderRepo.CreateOrderRoom(ctx, orderRoom); err != nil {
		if common.IsForeignKeyViolation(err) {
			return 0, "", common.ErrRoomNotFound
		}
		s.logger.Error("create order room failed", zap.Error(err))
		return 0, "", err
	}

	secretCode := common.GenerateBase58ID(16)
	orderData := types.OrderRoomData{
		ID:        id,
		ExpiredAt: booking.CheckOut,
	}
	bytes, _ := json.Marshal(orderData)

	redisKey := fmt.Sprintf("instay:order-room:%s", secretCode)
	ttl := booking.CheckOut.Sub(time.Now())

	if err = s.cacheProvider.SetObject(ctx, redisKey, bytes, ttl); err != nil {
		s.logger.Error("save order room data failed", zap.Error(err))
		return 0, "", err
	}

	return id, secretCode, nil
}

func (s *orderSvcImpl) VerifyOrderRoom(ctx context.Context, secretCode string) (string, time.Duration, error) {
	redisKey := fmt.Sprintf("instay:order-room:%s", secretCode)
	bytes, err := s.cacheProvider.GetObject(ctx, redisKey)
	if err != nil {
		s.logger.Error("get order room data failed", zap.Error(err))
		return "", 0, err
	}
	if bytes == nil {
		return "", 0, common.ErrInvalidToken
	}

	var orderRoomData types.OrderRoomData
	if err = json.Unmarshal(bytes, &orderRoomData); err != nil {
		s.logger.Error("unmarshal order room data failed", zap.Error(err))
		return "", 0, err
	}

	ttl := orderRoomData.ExpiredAt.Sub(time.Now())

	guestToken, err := s.jwtProvider.GenerateGuestToken(orderRoomData.ID, ttl)
	if err != nil {
		s.logger.Error("generate guest token failed", zap.Error(err))
		return "", 0, err
	}

	return guestToken, ttl, nil
}
