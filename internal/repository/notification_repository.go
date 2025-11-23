package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"gorm.io/gorm"
)

type Notification interface {
	CreateNotificationTx(ctx context.Context, tx *gorm.DB, notification *model.Notification) error
}