package service

import (
	"context"

	"github.com/MamangRust/monolith-pointofsale-pkg/kafka"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-pkg/outbox"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	mencache "github.com/MamangRust/monolith-pointofsale-transacton/cache"
	"github.com/MamangRust/monolith-pointofsale-transacton/repository"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	TransactionQuery           TransactionQueryService
	TransactionCommand         TransactionCommandService
	TransactionStats           TransactionStatsService
	TransactionStatsByMerchant TransactionStatsByMerchantService
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
	// Normalize a typed-nil *kafka.Kafka into a nil interface so the
	// graceful-degradation guard (s.kafka == nil) in the command service works.
	var kafkaPublisher EmailEventPublisher
	if deps.Kafka != nil {
		kafkaPublisher = deps.Kafka
	}

	return &Service{
		TransactionQuery:           NewTransactionQueryService(deps.Mencache, deps.Repositories.TransactionQueryRepository, deps.Logger, deps.Observability),
		TransactionCommand:         NewTransactionCommandService(kafkaPublisher, deps.Mencache, deps.Repositories.CashierQuery, deps.Repositories.MerchantQuery, deps.Repositories.TransactionQueryRepository, deps.Repositories.TransactionCommandRepository, deps.Repositories.OrderQuery, deps.Repositories.OrderItemQuery, deps.Pool, deps.Outbox, deps.Logger, deps.Observability),
		TransactionStats:           NewTransactionStatsService(deps.Mencache, deps.Repositories.TransactionStatsRepository, deps.Logger, deps.Observability),
		TransactionStatsByMerchant: NewTransactionStatsByMerchantService(deps.Mencache, deps.Repositories.TransactionStatsByMerchant, deps.Logger, deps.Observability),
	}
}
