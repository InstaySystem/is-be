package service

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

type AuthService interface {
	Login(ctx context.Context, req types.LoginRequest) (*model.User, string, string, error)
}