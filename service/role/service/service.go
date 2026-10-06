package service

import (
	"context"

	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	mencache "github.com/MamangRust/monolith-pointofsale-role/cache"
	"github.com/MamangRust/monolith-pointofsale-role/repository"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
)

type Service struct {
	RoleQuery   RoleQueryService
	RoleCommand RoleCommandService
}

type Deps struct {
	Ctx           context.Context
	Mencache      mencache.Mencache
	Repositories  *repository.Repositories
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	return &Service{
		RoleQuery:   NewRoleQueryService(deps.Mencache, deps.Repositories.RoleQuery, deps.Logger, deps.Observability),
		RoleCommand: NewRoleCommandService(deps.Mencache, deps.Repositories.RoleCommand, deps.Logger, deps.Observability),
	}
}
