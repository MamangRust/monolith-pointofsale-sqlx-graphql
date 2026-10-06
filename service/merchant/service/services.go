package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-pointofsale-merchant/cache"
	"github.com/MamangRust/monolith-pointofsale-merchant/repository"
	"github.com/MamangRust/monolith-pointofsale-pkg/kafka"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-pkg/outbox"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	MerchantQuery           MerchantQueryService
	MerchantCommand         MerchantCommandService
	MerchantDocumentCommand MerchantDocumentCommandService
	MerchantDocumentQuery   MerchantDocumentQueryService
}

type Deps struct {
	Ctx           context.Context
	Kafka         *kafka.Kafka
	Mencache      mencache.Mencache
	Repositories  *repository.Repositories
	Pool          *sqlx.DB
	Outbox        *outbox.OutboxService
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	return &Service{
		MerchantQuery: NewMerchantQueryService(&merchantQueryDeps{
			Cache:         deps.Mencache,
			MerchantQuery: deps.Repositories.MerchantQuery,
			Logger:        deps.Logger,
			Observability: deps.Observability,
		}),
		MerchantCommand: NewMerchantCommandService(&merchantCommandDeps{
			Kafka:           deps.Kafka,
			Cache:           deps.Mencache,
			UserQuery:       deps.Repositories.UserQuery,
			MerchantQuery:   deps.Repositories.MerchantQuery,
			MerchantCommand: deps.Repositories.MerchantCommand,
			Pool:            deps.Pool,
			Outbox:          deps.Outbox,
			Logger:          deps.Logger,
			Observability:   deps.Observability,
		}),
		MerchantDocumentCommand: NewMerchantDocumentCommandService(&merchantDocumentCommandDeps{
			Kafka:                   deps.Kafka,
			Cache:                   deps.Mencache,
			MerchantQuery:           deps.Repositories.MerchantQuery,
			MerchantDocumentCommand: deps.Repositories.MerchantDocumentCommand,
			UserQuery:               deps.Repositories.UserQuery,
			Pool:                    deps.Pool,
			Outbox:                  deps.Outbox,
			Logger:                  deps.Logger,
			Observability:           deps.Observability,
		}),
		MerchantDocumentQuery: NewMerchantDocumentQueryService(&merchantDocumentQueryDeps{
			Cache:                 deps.Mencache,
			MerchantDocumentQuery: deps.Repositories.MerchantDocumentQuery,
			Logger:                deps.Logger,
			Observability:         deps.Observability,
		}),
	}
}
