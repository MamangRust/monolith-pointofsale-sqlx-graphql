package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	mencache "github.com/MamangRust/monolith-pointofsale-auth/cache"
	"github.com/MamangRust/monolith-pointofsale-auth/repository"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/email"
	"github.com/MamangRust/monolith-pointofsale-pkg/event"
	"github.com/MamangRust/monolith-pointofsale-pkg/hash"
	"github.com/MamangRust/monolith-pointofsale-pkg/kafka"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-pkg/outbox"
	"github.com/MamangRust/monolith-pointofsale-pkg/randomstring"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-pointofsale-shared/errorhandler"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/user_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type RegisterServiceDeps struct {
	Cache mencache.RegisterCache

	User repository.UserRepository

	Role repository.RoleRepository

	UserRole repository.UserRoleRepository

	Hash hash.HashPassword

	Kafka *kafka.Kafka

	Outbox *outbox.OutboxService

	Logger logger.LoggerInterface

	Observability observability.TraceLoggerObservability
}

type registerService struct {
	mencache mencache.RegisterCache

	user repository.UserRepository

	role repository.RoleRepository

	userRole repository.UserRoleRepository

	hash hash.HashPassword

	kafka *kafka.Kafka

	outbox *outbox.OutboxService

	logger logger.LoggerInterface

	observability observability.TraceLoggerObservability
}

func NewRegisterService(params *RegisterServiceDeps) *registerService {
	return &registerService{
		mencache:      params.Cache,
		user:          params.User,
		role:          params.Role,
		userRole:      params.UserRole,
		hash:          params.Hash,
		kafka:         params.Kafka,
		outbox:        params.Outbox,
		logger:        params.Logger,
		observability: params.Observability,
	}
}

func (s *registerService) Register(ctx context.Context, request *requests.RegisterRequest) (*db.User, error) {
	const method = "Register"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("email", request.Email))

	defer func() {
		end(status)
	}()

	existingUser, err := s.user.FindByEmail(ctx, request.Email)
	if err == nil && existingUser != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](
			s.logger,
			user_errors.ErrUserEmailAlready,
			method,
			span,
			zap.String("email", request.Email),
		)
	}

	passwordHash, err := s.hash.HashPassword(request.Password)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](s.logger, err, method, span)
	}
	request.Password = passwordHash

	const defaultRoleName = "ROLE_ADMIN"
	role, err := s.role.FindByName(ctx, defaultRoleName)
	if err != nil || role == nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](s.logger, err, method, span, zap.String("role_name", defaultRoleName))
	}

	random, err := randomstring.GenerateRandomString(10)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](s.logger, err, method, span)
	}
	request.VerifiedCode = random
	// Auto-verify in non-production environments so the full register → verify →
	// login journey is testable locally; production keeps email verification.
	switch strings.ToLower(viper.GetString("APP_ENV")) {
	case "production", "kubernetes":
		request.IsVerified = false
	default:
		request.IsVerified = true
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Welcome to SanEdge",
		"Message": "Your account has been successfully created.",
		"Button":  "Verify Now",
		"Link":    "https://sanedge.example.com/login?verify_code=" + request.VerifiedCode,
	})

	payloadBytes, err := event.MarshalEmail("auth.register", request.Email, "Welcome to SanEdge", htmlBody)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](s.logger, err, method, span)
	}

	// User and role writes are owned by the user/role services and reached over
	// gRPC, so they cannot share a local Postgres transaction. The user is
	// created first, then the role is assigned; a failure leaves the user
	// without a role, which is safe to retry.
	newUser, err := s.user.CreateUser(ctx, request)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](s.logger, err, method, span)
	}

	_, err = s.userRole.AssignRoleToUser(ctx, &requests.CreateUserRoleRequest{
		UserId: int(newUser.UserID),
		RoleId: int(role.RoleID),
	})
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](s.logger, err, method, span, zap.Int("user.id", int(newUser.UserID)))
	}

	if s.outbox != nil {
		if err := s.outbox.Enqueue(ctx, "email-service-topic-auth-register", strconv.Itoa(int(newUser.UserID)), payloadBytes); err != nil {
			s.logger.Error("failed to enqueue registration email", zap.Error(err), zap.String("email", request.Email))
		}
	} else if s.kafka != nil {
		if sendErr := s.kafka.SendMessage(ctx, "email-service-topic-auth-register", strconv.Itoa(int(newUser.UserID)), payloadBytes); sendErr != nil {
			s.logger.Error("failed to send registration email via kafka", zap.Error(sendErr), zap.String("email", request.Email))
		}
	}

	s.mencache.SetVerificationCodeCache(ctx, request.Email, random, 15*time.Minute)

	logSuccess("User registered successfully",
		zap.String("email", request.Email),
		zap.String("first_name", request.FirstName),
		zap.String("last_name", request.LastName),
	)

	return newUser, nil
}
