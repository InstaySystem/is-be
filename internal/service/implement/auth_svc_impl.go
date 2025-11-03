package implement

import (
	"context"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/config"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/provider/jwt"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/InstaySystem/is-be/pkg/bcrypt"
	"go.uber.org/zap"
)

type authSvcImpl struct {
	userRepo    repository.UserRepository
	logger      *zap.Logger
	bHash       bcrypt.Hasher
	jwtProvider jwt.JWTProvider
	cfg         *config.Config
}

func NewAuthService(userRepo repository.UserRepository, logger *zap.Logger, bHash bcrypt.Hasher, jwtProvider jwt.JWTProvider, cfg *config.Config) service.AuthService {
	return &authSvcImpl{
		userRepo,
		logger,
		bHash,
		jwtProvider,
		cfg,
	}
}

func (s *authSvcImpl) Login(ctx context.Context, req types.LoginRequest) (*model.User, string, string, error) {
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		s.logger.Error("find user by username failed", zap.String("username", req.Username), zap.Error(err))
		return nil, "", "", err
	}
	if user == nil {
		return nil, "", "", common.ErrLoginFailed
	}

	if err = s.bHash.VerifyPassword(req.Password, user.Password); err != nil {
		return nil, "", "", common.ErrLoginFailed
	}

	accessToken, err := s.jwtProvider.GenerateToken(user.ID, user.Role, s.cfg.JWT.AccessExpiresIn)
	if err != nil {
		s.logger.Error("generate access token failed", zap.Error(err))
		return nil, "", "", err
	}

	refreshToken, err := s.jwtProvider.GenerateToken(user.ID, user.Role, s.cfg.JWT.RefreshExpiresIn)
	if err != nil {
		s.logger.Error("generate refresh token failed", zap.Error(err))
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}
