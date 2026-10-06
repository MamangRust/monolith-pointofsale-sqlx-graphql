package service

import (
	"context"
	"fmt"
	"strconv"

	mencache "github.com/MamangRust/monolith-pointofsale-merchant/cache"
	"github.com/MamangRust/monolith-pointofsale-merchant/repository"
	useradapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/user"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/email"
	"github.com/MamangRust/monolith-pointofsale-pkg/event"
	"github.com/MamangRust/monolith-pointofsale-pkg/kafka"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-pkg/outbox"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-pointofsale-shared/errorhandler"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/merchant_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/user_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantCommandDeps struct {
	Kafka           *kafka.Kafka
	Cache           mencache.MerchantCommandCache
	UserQuery       useradapter.QueryRepository
	MerchantQuery   repository.MerchantQueryRepository
	MerchantCommand repository.MerchantCommandRepository
	Pool            *sqlx.DB
	Outbox          *outbox.OutboxService
	Logger          logger.LoggerInterface
	Observability   observability.TraceLoggerObservability
}

type merchantCommandService struct {
	kafka                     *kafka.Kafka
	mencache                  mencache.MerchantCommandCache
	userRepository            useradapter.QueryRepository
	merchantQueryRepository   repository.MerchantQueryRepository
	merchantCommandRepository repository.MerchantCommandRepository
	pool                      *sqlx.DB
	outbox                    *outbox.OutboxService
	logger                    logger.LoggerInterface
	observability             observability.TraceLoggerObservability
}

func NewMerchantCommandService(params *merchantCommandDeps) MerchantCommandService {
	return &merchantCommandService{
		kafka:                     params.Kafka,
		mencache:                  params.Cache,
		userRepository:            params.UserQuery,
		merchantQueryRepository:   params.MerchantQuery,
		merchantCommandRepository: params.MerchantCommand,
		pool:                      params.Pool,
		outbox:                    params.Outbox,
		logger:                    params.Logger,
		observability:             params.Observability,
	}
}

func (s *merchantCommandService) CreateMerchant(ctx context.Context, request *requests.CreateMerchantRequest) (*db.Merchant, error) {
	const method = "CreateMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user.id", request.UserID),
	)
	defer func() {
		end(status)
	}()

	user, err := s.userRepository.FindById(ctx, request.UserID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Merchant](
			s.logger,
			user_errors.ErrUserNotFoundRes.WithInternal(err),
			method,
			span,
			zap.Int("user.id", request.UserID),
		)
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Welcome to SanEdge Merchant Portal",
		"Message": "Your merchant account has been created successfully. To continue, please upload the required documents for verification. Once completed, our team will review and activate your account.",
		"Button":  "Upload Documents",
		"Link":    fmt.Sprintf("https://sanedge.example.com/merchant/%d/documents", user.UserID),
	})

	payloadBytes, err := event.MarshalEmail("merchant.created", user.Email, "Initial Verification - SanEdge", htmlBody)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Merchant](
			s.logger,
			merchant_errors.ErrFailedSendEmail.WithInternal(err),
			method,
			span,
			zap.Int("user.id", int(user.UserID)),
		)
	}

	var res *db.Merchant
	if s.pool != nil {
		tx, beginErr := s.pool.Beginx()
		if beginErr != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.Merchant](s.logger, beginErr, method, span, zap.Int("user.id", request.UserID))
		}
		defer func() { _ = tx.Rollback() }()

		res, err = s.merchantCommandRepository.CreateMerchantInTx(ctx, tx, request)
		if err == nil && s.outbox != nil {
			err = s.outbox.EnqueueInTx(ctx, tx, "email-service-topic-merchant-create", strconv.Itoa(int(res.MerchantID)), payloadBytes)
		}
		if err != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.Merchant](
				s.logger,
				merchant_errors.ErrFailedCreateMerchant.WithInternal(err),
				method,
				span,
				zap.Int("user.id", request.UserID),
			)
		}
		if err := tx.Commit(); err != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.Merchant](
				s.logger,
				merchant_errors.ErrFailedCreateMerchant.WithInternal(err),
				method,
				span,
				zap.Int("user.id", request.UserID),
			)
		}
	} else {
		res, err = s.merchantCommandRepository.CreateMerchant(ctx, request)
		if err != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.Merchant](
				s.logger,
				merchant_errors.ErrFailedCreateMerchant.WithInternal(err),
				method,
				span,
				zap.Int("user.id", request.UserID),
			)
		}
		if s.kafka != nil {
			if sendErr := s.kafka.SendMessage(ctx, "email-service-topic-merchant-create", strconv.Itoa(int(res.MerchantID)), payloadBytes); sendErr != nil {
				// Graceful degradation: email failure must not fail merchant creation.
				s.logger.Warn("failed to send merchant created email via kafka",
					zap.Int("merchant.id", int(res.MerchantID)),
					zap.Error(sendErr),
				)
			}
		}
	}

	logSuccess("Successfully created merchant", zap.Int("merchant.id", int(res.MerchantID)))
	return res, nil
}

func (s *merchantCommandService) UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*db.Merchant, error) {
	const method = "UpdateMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant.id", *request.MerchantID),
	)
	defer func() {
		end(status)
	}()

	res, err := s.merchantCommandRepository.UpdateMerchant(ctx, request)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Merchant](
			s.logger,
			merchant_errors.ErrFailedUpdateMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteCachedMerchant(ctx, *request.MerchantID)
	logSuccess("Successfully updated merchant", zap.Int("merchant.id", *request.MerchantID))
	return res, nil
}

func (s *merchantCommandService) UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*db.Merchant, error) {
	const method = "UpdateMerchantStatus"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant.id", *request.MerchantID),
	)
	defer func() {
		end(status)
	}()

	merchant, err := s.merchantQueryRepository.FindById(ctx, *request.MerchantID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Merchant](
			s.logger,
			merchant_errors.ErrFailedFindMerchantById.WithInternal(err),
			method,
			span,
			zap.Int("merchant.id", *request.MerchantID),
		)
	}

	user, err := s.userRepository.FindById(ctx, int(merchant.UserID))
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Merchant](
			s.logger,
			user_errors.ErrUserNotFoundRes.WithInternal(err),
			method,
			span,
			zap.Int("user.id", int(merchant.UserID)),
		)
	}

	statusReq := request.Status
	subject := ""
	message := ""
	buttonLabel := "Go to Portal"
	link := fmt.Sprintf("https://sanedge.example.com/merchant/%d/dashboard", *request.MerchantID)

	switch statusReq {
	case "active":
		subject = "Your Merchant Account is Now Active"
		message = "Congratulations! Your merchant account has been verified and is now <b>active</b>. You can now fully access all features in the SanEdge Merchant Portal."
	case "inactive":
		subject = "Merchant Account Set to Inactive"
		message = "Your merchant account status has been set to <b>inactive</b>. Please contact support if you believe this is a mistake."
	case "rejected":
		subject = "Merchant Account Rejected"
		message = "We're sorry to inform you that your merchant account has been <b>rejected</b>. Please contact support or review your submissions."
	default:
		return nil, nil
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   subject,
		"Message": message,
		"Button":  buttonLabel,
		"Link":    link,
	})

	payloadBytes, err := event.MarshalEmail("merchant.status_updated", user.Email, subject, htmlBody)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Merchant](
			s.logger,
			merchant_errors.ErrFailedSendEmail.WithInternal(err),
			method,
			span,
			zap.Int("merchant.id", *request.MerchantID),
		)
	}

	var res *db.Merchant
	if s.pool != nil {
		tx, beginErr := s.pool.Beginx()
		if beginErr != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.Merchant](s.logger, beginErr, method, span, zap.Int("merchant.id", *request.MerchantID))
		}
		defer func() { _ = tx.Rollback() }()

		res, err = s.merchantCommandRepository.UpdateMerchantStatusInTx(ctx, tx, request)
		if err == nil && s.outbox != nil {
			err = s.outbox.EnqueueInTx(ctx, tx, "email-service-topic-merchant-update-status", strconv.Itoa(*request.MerchantID), payloadBytes)
		}
		if err != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.Merchant](
				s.logger,
				merchant_errors.ErrFailedUpdateMerchant.WithInternal(err),
				method,
				span,
				zap.Int("merchant.id", *request.MerchantID),
			)
		}
		if err := tx.Commit(); err != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.Merchant](
				s.logger,
				merchant_errors.ErrFailedUpdateMerchant.WithInternal(err),
				method,
				span,
				zap.Int("merchant.id", *request.MerchantID),
			)
		}
	} else {
		res, err = s.merchantCommandRepository.UpdateMerchantStatus(ctx, request)
		if err != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.Merchant](
				s.logger,
				merchant_errors.ErrFailedUpdateMerchant.WithInternal(err),
				method,
				span,
				zap.Int("merchant.id", *request.MerchantID),
			)
		}
		if s.kafka != nil {
			if sendErr := s.kafka.SendMessage(ctx, "email-service-topic-merchant-update-status", strconv.Itoa(*request.MerchantID), payloadBytes); sendErr != nil {
				// Graceful degradation: email failure must not fail merchant status update.
				s.logger.Warn("failed to send merchant status email via kafka",
					zap.Int("merchant.id", *request.MerchantID),
					zap.Error(sendErr),
				)
			}
		}
	}

	s.mencache.DeleteCachedMerchant(ctx, *request.MerchantID)
	logSuccess("Successfully updated merchant status", zap.Int("merchant.id", *request.MerchantID))
	return res, nil
}

func (s *merchantCommandService) TrashedMerchant(ctx context.Context, merchant_id int) (*db.Merchant, error) {
	const method = "TrashedMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant.id", merchant_id),
	)
	defer func() {
		end(status)
	}()

	res, err := s.merchantCommandRepository.TrashedMerchant(ctx, merchant_id)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Merchant](
			s.logger,
			merchant_errors.ErrFailedTrashMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteCachedMerchant(ctx, merchant_id)
	logSuccess("Successfully trashed merchant", zap.Int("merchant.id", merchant_id))
	return res, nil
}

func (s *merchantCommandService) RestoreMerchant(ctx context.Context, merchant_id int) (*db.Merchant, error) {
	const method = "RestoreMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant.id", merchant_id),
	)
	defer func() {
		end(status)
	}()

	res, err := s.merchantCommandRepository.RestoreMerchant(ctx, merchant_id)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Merchant](
			s.logger,
			merchant_errors.ErrFailedRestoreMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteCachedMerchant(ctx, merchant_id)
	logSuccess("Successfully restored merchant", zap.Int("merchant.id", merchant_id))
	return res, nil
}

func (s *merchantCommandService) DeleteMerchantPermanent(ctx context.Context, merchant_id int) (bool, error) {
	const method = "DeleteMerchantPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant.id", merchant_id),
	)
	defer func() {
		end(status)
	}()

	success, err := s.merchantCommandRepository.DeleteMerchantPermanent(ctx, merchant_id)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			merchant_errors.ErrFailedDeleteMerchantPermanent.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteCachedMerchant(ctx, merchant_id)
	logSuccess("Successfully deleted merchant permanently", zap.Int("merchant.id", merchant_id), zap.Bool("success", success))
	return success, nil
}

func (s *merchantCommandService) RestoreAllMerchant(ctx context.Context) (bool, error) {
	const method = "RestoreAllMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	success, err := s.merchantCommandRepository.RestoreAllMerchant(ctx)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			merchant_errors.ErrFailedRestoreAllMerchants.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteCachedMerchantAllCache(ctx)
	logSuccess("Successfully restored all merchants", zap.Bool("success", success))
	return success, nil
}

func (s *merchantCommandService) DeleteAllMerchantPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllMerchantPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	success, err := s.merchantCommandRepository.DeleteAllMerchantPermanent(ctx)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			merchant_errors.ErrFailedDeleteAllMerchantsPermanent.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteCachedMerchantAllCache(ctx)
	logSuccess("Successfully deleted all merchants permanently", zap.Bool("success", success))
	return success, nil
}
