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
	merchantdocument_errors "github.com/MamangRust/monolith-pointofsale-shared/errors/merchant_document_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/merchant_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/user_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantDocumentCommandDeps struct {
	Kafka                   *kafka.Kafka
	Cache                   mencache.MerchantDocumentCommandCache
	MerchantQuery           repository.MerchantQueryRepository
	MerchantDocumentCommand repository.MerchantDocumentCommandRepository
	UserQuery               useradapter.QueryRepository
	Pool                    *sqlx.DB
	Outbox                  *outbox.OutboxService
	Logger                  logger.LoggerInterface
	Observability           observability.TraceLoggerObservability
}

type merchantDocumentCommandService struct {
	kafka                             *kafka.Kafka
	mencache                          mencache.MerchantDocumentCommandCache
	merchantQueryRepository           repository.MerchantQueryRepository
	merchantDocumentCommandRepository repository.MerchantDocumentCommandRepository
	userRepository                    useradapter.QueryRepository
	pool                              *sqlx.DB
	outbox                            *outbox.OutboxService
	logger                            logger.LoggerInterface
	observability                     observability.TraceLoggerObservability
}

func NewMerchantDocumentCommandService(params *merchantDocumentCommandDeps) MerchantDocumentCommandService {
	return &merchantDocumentCommandService{
		kafka:                             params.Kafka,
		mencache:                          params.Cache,
		merchantQueryRepository:           params.MerchantQuery,
		merchantDocumentCommandRepository: params.MerchantDocumentCommand,
		userRepository:                    params.UserQuery,
		pool:                              params.Pool,
		outbox:                            params.Outbox,
		logger:                            params.Logger,
		observability:                     params.Observability,
	}
}

func (s *merchantDocumentCommandService) CreateMerchantDocument(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*db.MerchantDocument, error) {
	const method = "CreateMerchantDocument"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant.id", request.MerchantID),
	)
	defer func() {
		end(status)
	}()

	merchant, err := s.merchantQueryRepository.FindById(ctx, request.MerchantID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.MerchantDocument](
			s.logger,
			merchant_errors.ErrFailedFindMerchantById,
			method,
			span,
			zap.Int("merchant.id", request.MerchantID),
		)
	}

	user, err := s.userRepository.FindById(ctx, int(merchant.UserID))
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.MerchantDocument](
			s.logger,
			user_errors.ErrUserNotFoundRes,
			method,
			span,
			zap.Int("user.id", int(merchant.UserID)),
		)
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Welcome to SanEdge Merchant Portal",
		"Message": "Thank you for registering your merchant account. Your account is currently <b>inactive</b> and under initial review. To proceed, please upload all required documents for verification. Once your documents are submitted, our team will review them and activate your account accordingly.",
		"Button":  "Upload Documents",
		"Link":    fmt.Sprintf("https://sanedge.example.com/merchant/%d/documents", user.UserID),
	})

	payloadBytes, err := event.MarshalEmail("merchant_document.created", user.Email, "Merchant Verification Pending - Action Required", htmlBody)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.MerchantDocument](
			s.logger,
			merchant_errors.ErrFailedSendEmail,
			method,
			span,
			zap.Int("merchant.id", request.MerchantID),
		)
	}

	var merchantDocument *db.MerchantDocument
	if s.pool != nil {
		tx, beginErr := s.pool.Beginx()
		if beginErr != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.MerchantDocument](s.logger, beginErr, method, span, zap.Int("merchant.id", request.MerchantID))
		}
		defer func() { _ = tx.Rollback() }()

		merchantDocument, err = s.merchantDocumentCommandRepository.CreateMerchantDocumentInTx(ctx, tx, request)
		if err == nil && s.outbox != nil {
			err = s.outbox.EnqueueInTx(ctx, tx, "email-service-topic-merchant-document-create", strconv.Itoa(int(merchantDocument.DocumentID)), payloadBytes)
		}
		if err != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.MerchantDocument](
				s.logger,
				merchantdocument_errors.ErrCreateMerchantDocumentFailed.WithInternal(err),
				method,
				span,
				zap.Int("merchant.id", request.MerchantID),
			)
		}
		if err := tx.Commit(); err != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.MerchantDocument](
				s.logger,
				merchantdocument_errors.ErrCreateMerchantDocumentFailed.WithInternal(err),
				method,
				span,
				zap.Int("merchant.id", request.MerchantID),
			)
		}
	} else {
		merchantDocument, err = s.merchantDocumentCommandRepository.CreateMerchantDocument(ctx, request)
		if err != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.MerchantDocument](
				s.logger,
				merchantdocument_errors.ErrCreateMerchantDocumentFailed,
				method,
				span,
				zap.Int("merchant.id", request.MerchantID),
			)
		}
		if s.kafka != nil {
			if sendErr := s.kafka.SendMessage(ctx, "email-service-topic-merchant-document-create", strconv.Itoa(int(merchantDocument.DocumentID)), payloadBytes); sendErr != nil {
				// Graceful degradation: email failure must not fail document creation.
				s.logger.Warn("failed to send merchant document created email via kafka",
					zap.Int("document.id", int(merchantDocument.DocumentID)),
					zap.Error(sendErr),
				)
			}
		}
	}

	logSuccess("Successfully created merchant document", zap.Int("merchantDocument.id", int(merchantDocument.DocumentID)))
	return merchantDocument, nil
}

func (s *merchantDocumentCommandService) UpdateMerchantDocument(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*db.MerchantDocument, error) {
	const method = "UpdateMerchantDocument"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("document_id", *request.DocumentID),
	)
	defer func() {
		end(status)
	}()

	merchantDocument, err := s.merchantDocumentCommandRepository.UpdateMerchantDocument(ctx, request)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.MerchantDocument](
			s.logger,
			merchantdocument_errors.ErrUpdateMerchantDocumentFailed,
			method,
			span,
			zap.Int("merchantDocument.id", *request.DocumentID),
		)
	}

	s.mencache.DeleteCachedMerchantDocuments(ctx, int(merchantDocument.DocumentID))
	logSuccess("Successfully updated merchant document", zap.Int("merchantDocument.id", *request.DocumentID))
	return merchantDocument, nil
}

func (s *merchantDocumentCommandService) UpdateMerchantDocumentStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*db.MerchantDocument, error) {
	const method = "UpdateMerchantDocumentStatus"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantDocument.id", *request.DocumentID),
	)
	defer func() {
		end(status)
	}()

	merchant, err := s.merchantQueryRepository.FindById(ctx, request.MerchantID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.MerchantDocument](
			s.logger,
			merchant_errors.ErrFailedFindMerchantById,
			method,
			span,
			zap.Int("merchant.id", request.MerchantID),
		)
	}

	user, err := s.userRepository.FindById(ctx, int(merchant.UserID))
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.MerchantDocument](
			s.logger,
			user_errors.ErrUserNotFoundRes,
			method,
			span,
			zap.Int("user.id", int(merchant.UserID)),
		)
	}

	statusReq := request.Status
	note := request.Note
	subject := ""
	message := ""
	buttonLabel := ""
	link := fmt.Sprintf("https://sanedge.example.com/merchant/%d/documents", request.MerchantID)
	sendEmail := true

	switch statusReq {
	case "pending":
		subject = "Merchant Document Status: Pending Review"
		message = "Your merchant documents have been submitted and are currently pending review."
		buttonLabel = "View Documents"
	case "approved":
		subject = "Merchant Document Status: Approved"
		message = "Congratulations! Your merchant documents have been approved. Your account is now active and fully functional."
		buttonLabel = "Go to Dashboard"
		link = fmt.Sprintf("https://sanedge.example.com/merchant/%d/dashboard", request.MerchantID)
	case "rejected":
		subject = "Merchant Document Status: Rejected"
		message = "Unfortunately, your merchant documents were rejected. Please review the feedback below and re-upload the necessary documents."
		buttonLabel = "Re-upload Documents"
	default:
		// Status tanpa template email (mis. "verified"): update di DB tetap
		// sukses, hanya notifikasi email yang di-skip. Jangan return nil —
		// dokumen harus dikembalikan dan cache di-invalidasi (return nil
		// setelah write sukses memicu nil-pointer di response mapper).
		sendEmail = false
	}

	var merchantDocument *db.MerchantDocument
	if s.pool != nil {
		tx, beginErr := s.pool.Beginx()
		if beginErr != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.MerchantDocument](s.logger, beginErr, method, span, zap.Int("merchantDocument.id", *request.DocumentID))
		}
		defer func() { _ = tx.Rollback() }()

		merchantDocument, err = s.merchantDocumentCommandRepository.UpdateMerchantDocumentStatusInTx(ctx, tx, request)
		if err == nil && sendEmail && s.outbox != nil {
			if note != "" {
				message += fmt.Sprintf(`<br><br><b>Reviewer Note:</b><br><i>%s</i>`, note)
			}

			htmlBody := email.GenerateEmailHTML(map[string]string{
				"Title":   subject,
				"Message": message,
				"Button":  buttonLabel,
				"Link":    link,
			})

			var payloadBytes []byte
			payloadBytes, err = event.MarshalEmail("merchant_document.status_updated", user.Email, subject, htmlBody)
			if err == nil {
				err = s.outbox.EnqueueInTx(ctx, tx, "email-service-topic-merchant-document-update-status", strconv.Itoa(request.MerchantID), payloadBytes)
			}
		}
		if err != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.MerchantDocument](
				s.logger,
				merchantdocument_errors.ErrUpdateMerchantDocumentFailed.WithInternal(err),
				method,
				span,
				zap.Int("merchantDocument.id", *request.DocumentID),
			)
		}
		if err := tx.Commit(); err != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.MerchantDocument](
				s.logger,
				merchantdocument_errors.ErrUpdateMerchantDocumentFailed.WithInternal(err),
				method,
				span,
				zap.Int("merchantDocument.id", *request.DocumentID),
			)
		}
	} else {
		merchantDocument, err = s.merchantDocumentCommandRepository.UpdateMerchantDocumentStatus(ctx, request)
		if err != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.MerchantDocument](
				s.logger,
				merchantdocument_errors.ErrUpdateMerchantDocumentFailed,
				method,
				span,
				zap.Int("merchantDocument.id", *request.DocumentID),
			)
		}
		if sendEmail && s.kafka != nil {
			if note != "" {
				message += fmt.Sprintf(`<br><br><b>Reviewer Note:</b><br><i>%s</i>`, note)
			}

			htmlBody := email.GenerateEmailHTML(map[string]string{
				"Title":   subject,
				"Message": message,
				"Button":  buttonLabel,
				"Link":    link,
			})

			payloadBytes, marshalErr := event.MarshalEmail("merchant_document.status_updated", user.Email, subject, htmlBody)
			if marshalErr != nil {
				status = "error"
				return sharederrorhandler.HandleError[*db.MerchantDocument](
					s.logger,
					merchant_errors.ErrFailedSendEmail.WithInternal(marshalErr),
					method,
					span,
					zap.Int("merchant.id", request.MerchantID),
				)
			}
			if sendErr := s.kafka.SendMessage(ctx, "email-service-topic-merchant-document-update-status", strconv.Itoa(request.MerchantID), payloadBytes); sendErr != nil {
				// Graceful degradation: email failure must not fail document status update.
				s.logger.Warn("failed to send merchant document status email via kafka",
					zap.Int("merchant.id", request.MerchantID),
					zap.Error(sendErr),
				)
			}
		}
	}

	s.mencache.DeleteCachedMerchantDocuments(ctx, int(merchantDocument.DocumentID))
	logSuccess("Successfully updated merchant document status", zap.Int("merchantDocument.id", *request.DocumentID))
	return merchantDocument, nil
}

func (s *merchantDocumentCommandService) TrashedMerchantDocument(ctx context.Context, documentID int) (*db.MerchantDocument, error) {
	const method = "TrashedDocument"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantDocument.id", documentID),
	)
	defer func() {
		end(status)
	}()

	res, err := s.merchantDocumentCommandRepository.TrashedMerchantDocument(ctx, documentID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.MerchantDocument](
			s.logger,
			merchantdocument_errors.ErrTrashedMerchantDocumentFailed,
			method,
			span,
			zap.Int("document_id", documentID),
		)
	}

	s.mencache.DeleteCachedMerchantDocuments(ctx, documentID)
	logSuccess("Successfully trashed document", zap.Int("merchantDocument.id", documentID))
	return res, nil
}

func (s *merchantDocumentCommandService) RestoreMerchantDocument(ctx context.Context, documentID int) (*db.MerchantDocument, error) {
	const method = "RestoreDocument"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantDocument.id", documentID),
	)
	defer func() {
		end(status)
	}()

	res, err := s.merchantDocumentCommandRepository.RestoreMerchantDocument(ctx, documentID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.MerchantDocument](
			s.logger,
			merchantdocument_errors.ErrRestoreMerchantDocumentFailed,
			method,
			span,
			zap.Int("merchantDocument.id", documentID),
		)
	}

	s.mencache.DeleteCachedMerchantDocuments(ctx, documentID)
	logSuccess("Successfully restored document", zap.Int("merchantDocument.id", documentID))
	return res, nil
}

func (s *merchantDocumentCommandService) DeleteMerchantDocumentPermanent(ctx context.Context, documentID int) (bool, error) {
	const method = "DeleteDocumentPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantDocument.id", documentID),
	)
	defer func() {
		end(status)
	}()

	success, err := s.merchantDocumentCommandRepository.DeleteMerchantDocumentPermanent(ctx, documentID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			merchantdocument_errors.ErrDeleteMerchantDocumentPermanentFailed,
			method,
			span,
			zap.Int("merchantDocument.id", documentID),
		)
	}

	s.mencache.DeleteCachedMerchantDocuments(ctx, documentID)
	logSuccess("Successfully deleted document permanently", zap.Int("merchantDocument.id", documentID))
	return success, nil
}

func (s *merchantDocumentCommandService) RestoreAllMerchantDocument(ctx context.Context) (bool, error) {
	const method = "RestoreAllDocuments"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	success, err := s.merchantDocumentCommandRepository.RestoreAllMerchantDocument(ctx)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			merchantdocument_errors.ErrRestoreAllMerchantDocumentsFailed,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteCachedMerchantDocumentsAllCache(ctx)
	logSuccess("Successfully restored all documents", zap.Bool("success", success))
	return success, nil
}

func (s *merchantDocumentCommandService) DeleteAllMerchantDocumentPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllDocumentsPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	success, err := s.merchantDocumentCommandRepository.DeleteAllMerchantDocumentPermanent(ctx)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			merchantdocument_errors.ErrDeleteAllMerchantDocumentsPermanentFailed,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteCachedMerchantDocumentsAllCache(ctx)
	logSuccess("Successfully deleted all documents permanently", zap.Bool("success", success))
	return success, nil
}
