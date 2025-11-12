package container

import (
	"github.com/InstaySystem/is-be/internal/handler"
	repoImpl "github.com/InstaySystem/is-be/internal/repository/implement"
	svcImpl "github.com/InstaySystem/is-be/internal/service/implement"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ServiceContainer struct {
	Hdl *handler.ServiceHandler
}

func NewServiceContainer(
	db *gorm.DB,
	sfGen snowflake.Generator,
	logger *zap.Logger,
) *ServiceContainer {
	repo := repoImpl.NewServiceRepository(db)
	svc := svcImpl.NewServiceService(repo, sfGen, logger)
	hdl := handler.NewServiceHandler(svc)

	return &ServiceContainer{hdl}
}
