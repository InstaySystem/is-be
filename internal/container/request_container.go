package container

import (
	"github.com/InstaySystem/is-be/internal/handler"
	repoImpl "github.com/InstaySystem/is-be/internal/repository/implement"
	svcImpl "github.com/InstaySystem/is-be/internal/service/implement"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RequestContainer struct {
	Hdl *handler.RequestHandler
}

func NewRequestContainer(
	db *gorm.DB,
	sfGen snowflake.Generator,
	logger *zap.Logger,
) *RequestContainer {
	repo := repoImpl.NewRequestRepository(db)
	svc := svcImpl.NewRequestService(repo, sfGen, logger)
	hdl := handler.NewRequestHandler(svc)

	return &RequestContainer{hdl}
}
