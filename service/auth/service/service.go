package service

import (
	mencache "github.com/MamangRust/monolith-pointofsale-auth/cache"
	"github.com/MamangRust/monolith-pointofsale-auth/repository"
	"github.com/MamangRust/monolith-pointofsale-pkg/auth"
	"github.com/MamangRust/monolith-pointofsale-pkg/hash"
	"github.com/MamangRust/monolith-pointofsale-pkg/kafka"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-pkg/outbox"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	Login         LoginService
	Register      RegistrationService
	PasswordReset PasswordResetService
	Identify      IdentifyService
}

type Deps struct {
	Mencache      mencache.Mencache
	Repositories  *repository.Repositories
	Token         auth.TokenManager
	Hash          hash.HashPassword
	Logger        logger.LoggerInterface
	Kafka         *kafka.Kafka
	Pool          *sqlx.DB
	Outbox        *outbox.OutboxService
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	tokenService := NewTokenService(deps.Repositories.RefreshToken, deps.Token, deps.Logger, deps.Observability)

	return &Service{
		Login: NewLoginService(&LoginServiceDeps{
			Cache:          deps.Mencache,
			Logger:         deps.Logger,
			Hash:           deps.Hash,
			UserRepository: deps.Repositories.User,
			RefreshToken:   deps.Repositories.RefreshToken,
			Token:          deps.Token,
			TokenService:   tokenService,
			Observability:  deps.Observability,
		}),
		Register: NewRegisterService(&RegisterServiceDeps{
			Cache:         deps.Mencache,
			User:          deps.Repositories.User,
			Role:          deps.Repositories.Role,
			UserRole:      deps.Repositories.UserRole,
			Hash:          deps.Hash,
			Kafka:         deps.Kafka,
			Outbox:        deps.Outbox,
			Logger:        deps.Logger,
			Observability: deps.Observability,
		}),
		PasswordReset: NewPasswordResetService(&PasswordResetServiceDeps{
			Cache:         deps.Mencache,
			Kafka:         deps.Kafka,
			Logger:        deps.Logger,
			Hash:          deps.Hash,
			User:          deps.Repositories.User,
			ResetToken:    deps.Repositories.ResetToken,
			Pool:          deps.Pool,
			Outbox:        deps.Outbox,
			Observability: deps.Observability,
		}),
		Identify: NewIdentityService(&IdentityServiceDeps{
			Cache:         deps.Mencache,
			Token:         deps.Token,
			RefreshToken:  deps.Repositories.RefreshToken,
			User:          deps.Repositories.User,
			Logger:        deps.Logger,
			TokenService:  tokenService,
			Observability: deps.Observability,
		}),
	}
}
