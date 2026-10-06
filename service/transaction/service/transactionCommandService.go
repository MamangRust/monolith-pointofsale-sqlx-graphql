package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	cashieradapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/cashier"
	merchantadapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/merchant"
	orderadapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/order"
	orderitemadapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/orderitem"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/email"
	"github.com/MamangRust/monolith-pointofsale-pkg/event"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-pkg/outbox"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-pointofsale-shared/errorhandler"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/cashier_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/merchant_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/order_errors"
	orderitem_errors "github.com/MamangRust/monolith-pointofsale-shared/errors/order_item_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/transaction_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	mencache "github.com/MamangRust/monolith-pointofsale-transacton/cache"
	"github.com/MamangRust/monolith-pointofsale-transacton/repository"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transactionCommandService struct {
	kafka                        EmailEventPublisher
	mencache                     mencache.TransactionCommandCache
	cashierQueryRepository       cashieradapter.QueryRepository
	merchantQueryRepository      merchantadapter.Repository
	transactionQueryRepository   repository.TransactionQueryRepository
	transactionCommandRepository repository.TransactionCommandRepository
	orderQueryRepository         orderadapter.QueryRepository
	orderItemQueryRepository     orderitemadapter.QueryRepository
	pool                         *sqlx.DB
	outbox                       *outbox.OutboxService
	logger                       logger.LoggerInterface
	observability                observability.TraceLoggerObservability
}

func NewTransactionCommandService(
	kafka EmailEventPublisher,
	mencache mencache.TransactionCommandCache,
	cashierQueryRepository cashieradapter.QueryRepository,
	merchantQueryRepository merchantadapter.Repository,
	transactionQueryRepository repository.TransactionQueryRepository,
	transactionCommandRepository repository.TransactionCommandRepository,
	orderQueryRepository orderadapter.QueryRepository,
	orderItemQueryRepository orderitemadapter.QueryRepository,
	pool *sqlx.DB,
	outbox *outbox.OutboxService,
	logger logger.LoggerInterface,
	obs observability.TraceLoggerObservability,
) *transactionCommandService {
	return &transactionCommandService{
		kafka:                        kafka,
		mencache:                     mencache,
		cashierQueryRepository:       cashierQueryRepository,
		merchantQueryRepository:      merchantQueryRepository,
		transactionQueryRepository:   transactionQueryRepository,
		transactionCommandRepository: transactionCommandRepository,
		orderQueryRepository:         orderQueryRepository,
		orderItemQueryRepository:     orderItemQueryRepository,
		pool:                         pool,
		outbox:                       outbox,
		logger:                       logger,
		observability:                obs,
	}
}

func (s *transactionCommandService) CreateTransaction(ctx context.Context, req *requests.CreateTransactionRequest) (*db.Transaction, error) {
	const method = "CreateTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("cashier.id", req.CashierID), attribute.Int("merchant.id", req.MerchantID), attribute.Int("order.id", req.OrderID))

	defer func() {
		end(status)
	}()

	cashier, err := s.cashierQueryRepository.FindById(ctx, req.CashierID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			cashier_errors.ErrFailedFindCashierById.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	merchant, err := s.merchantQueryRepository.FindById(ctx, int(cashier.MerchantID))
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			merchant_errors.ErrFailedFindMerchantById.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}
	req.MerchantID = int(cashier.MerchantID)

	_, err = s.orderQueryRepository.FindById(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			order_errors.ErrFailedFindOrderById.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			orderitem_errors.ErrFailedFindOrderItemByOrder.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	if len(orderItems) == 0 {
		status = "error"
		errEmpty := errors.New("order item is empty")
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedOrderItemEmpty.WithInternal(errEmpty),
			method,
			span,
			zap.Error(errEmpty),
		)
	}

	var totalAmount int
	for _, item := range orderItems {
		if item.Quantity <= 0 {
			status = "error"
			errQty := errors.New("invalid quantity for order item")
			return sharederrorhandler.HandleError[*db.Transaction](
				s.logger,
				transaction_errors.ErrFailedPaymentStatusInvalid.WithInternal(errQty),
				method,
				span,
				zap.Error(errQty),
			)
		}
		totalAmount += int(item.Price) * int(item.Quantity)
	}

	ppn := totalAmount * 11 / 100
	totalAmountWithTax := totalAmount + ppn

	span.SetAttributes(
		attribute.Int("amount.subtotal", totalAmount),
		attribute.Int("amount.tax", ppn),
		attribute.Int("amount.total", totalAmountWithTax),
	)

	var paymentStatus string
	if req.Amount >= totalAmountWithTax {
		paymentStatus = "success"
	} else {
		status = "error"
		errBalance := errors.New("insufficient balance")
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedPaymentInsufficientBalance.WithInternal(errBalance),
			method,
			span,
			zap.Error(errBalance),
		)
	}

	req.Amount = totalAmountWithTax
	req.PaymentStatus = &paymentStatus

	var transaction *db.Transaction
	if s.pool != nil {
		tx, beginErr := s.pool.Beginx()
		if beginErr != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.Transaction](
				s.logger,
				transaction_errors.ErrFailedCreateTransaction.WithInternal(beginErr),
				method,
				span,
				zap.Error(beginErr),
			)
		}
		defer func() { _ = tx.Rollback() }()

		transaction, err = s.transactionCommandRepository.CreateTransactionInTx(ctx, tx, req)
		if err == nil && s.outbox != nil {
			err = s.enqueueTransactionCreateEvent(ctx, tx, merchant, transaction)
		}
		if err != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.Transaction](
				s.logger,
				transaction_errors.ErrFailedCreateTransaction.WithInternal(err),
				method,
				span,
				zap.Error(err),
			)
		}
		if err := tx.Commit(); err != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.Transaction](
				s.logger,
				transaction_errors.ErrFailedCreateTransaction.WithInternal(err),
				method,
				span,
				zap.Error(err),
			)
		}
	} else {
		transaction, err = s.transactionCommandRepository.CreateTransaction(ctx, req)
		if err != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.Transaction](
				s.logger,
				transaction_errors.ErrFailedCreateTransaction.WithInternal(err),
				method,
				span,
				zap.Error(err),
			)
		}

		// Fallback: fire-and-forget direct publish (tests/local only). The active
		// trace context is carried over so the consumer can continue the trace.
		eventCtx := trace.ContextWithSpanContext(
			context.WithoutCancel(ctx),
			trace.SpanContextFromContext(ctx),
		)
		go s.sendTransactionCreateEvent(eventCtx, merchant, transaction)
	}

	logSuccess("Successfully created transaction", zap.Bool("success", true))

	return transaction, nil
}

// enqueueTransactionCreateEvent builds the transaction email notification
// envelope and enqueues it in the given database transaction (Phase 6 —
// transactional outbox). Graceful degradation: when the merchant contact email
// is unavailable or marshaling fails, a warning is logged and the request stays
// successful (the event is simply not enqueued).
func (s *transactionCommandService) enqueueTransactionCreateEvent(ctx context.Context, tx *sqlx.Tx, merchant *db.Merchant, transaction *db.Transaction) error {
	if s.outbox == nil {
		return nil
	}
	if merchant == nil || merchant.ContactEmail == nil || *merchant.ContactEmail == "" {
		s.logger.Warn("Merchant contact email not available; skipping transaction email notification",
			zap.Int("transaction.id", int(transaction.TransactionID)),
			zap.Int("merchant.id", int(transaction.MerchantID)),
		)
		return nil
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Transaction Successful",
		"Message": fmt.Sprintf("Your transaction #%d of amount %d has been completed successfully.", transaction.TransactionID, transaction.Amount),
		"Button":  "View Transaction",
		"Link":    fmt.Sprintf("https://sanedge.example.com/transactions/%d", transaction.TransactionID),
	})

	payloadBytes, err := event.MarshalEmail("transaction.created", *merchant.ContactEmail, "Transaction Successful - SanEdge", htmlBody)
	if err != nil {
		s.logger.Warn("failed to marshal transaction email payload",
			zap.Int("transaction.id", int(transaction.TransactionID)),
			zap.Error(err),
		)
		return err
	}

	return s.outbox.EnqueueInTx(ctx, tx, "email-service-topic-transaction-create", strconv.Itoa(int(transaction.TransactionID)), payloadBytes)
}

// sendTransactionCreateEvent publishes the transaction email notification
// event to the `email-service-topic-transaction-create` topic (consumed by the
// email service). Graceful degradation: when Kafka is not configured, the
// merchant contact email is empty, marshaling fails, or publishing fails →
// log a warning and keep the request successful.
func (s *transactionCommandService) sendTransactionCreateEvent(ctx context.Context, merchant *db.Merchant, transaction *db.Transaction) {
	if s.kafka == nil {
		s.logger.Warn("Kafka not configured; skipping transaction email notification",
			zap.Int("transaction.id", int(transaction.TransactionID)),
			zap.Int("merchant.id", int(transaction.MerchantID)),
		)
		return
	}

	if merchant == nil || merchant.ContactEmail == nil || *merchant.ContactEmail == "" {
		s.logger.Warn("Merchant contact email not available; skipping transaction email notification",
			zap.Int("transaction.id", int(transaction.TransactionID)),
			zap.Int("merchant.id", int(transaction.MerchantID)),
		)
		return
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Transaction Successful",
		"Message": fmt.Sprintf("Your transaction #%d of amount %d has been completed successfully.", transaction.TransactionID, transaction.Amount),
		"Button":  "View Transaction",
		"Link":    fmt.Sprintf("https://sanedge.example.com/transactions/%d", transaction.TransactionID),
	})

	payloadBytes, err := event.MarshalEmail("transaction.created", *merchant.ContactEmail, "Transaction Successful - SanEdge", htmlBody)
	if err != nil {
		s.logger.Warn("failed to marshal transaction email payload",
			zap.Int("transaction.id", int(transaction.TransactionID)),
			zap.Error(err),
		)
		return
	}

	if err := s.kafka.SendMessage(ctx, "email-service-topic-transaction-create", strconv.Itoa(int(transaction.TransactionID)), payloadBytes); err != nil {
		s.logger.Warn("failed to send transaction email via kafka",
			zap.Int("transaction.id", int(transaction.TransactionID)),
			zap.Error(err),
		)
	}
}

func (s *transactionCommandService) UpdateTransaction(ctx context.Context, req *requests.UpdateTransactionRequest) (*db.Transaction, error) {
	const method = "UpdateTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transaction.id", *req.TransactionID), attribute.Int("merchant.id", req.MerchantID), attribute.Int("order.id", req.OrderID))

	defer func() {
		end(status)
	}()

	cashier, err := s.cashierQueryRepository.FindById(ctx, req.CashierID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			cashier_errors.ErrFailedFindCashierById.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	existingTx, err := s.transactionQueryRepository.FindById(ctx, *req.TransactionID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedFindTransactionById.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	if existingTx.PaymentStatus == "paid" || existingTx.PaymentStatus == "refunded" {
		status = "error"
		errStatus := errors.New("transaction payment status cannot be modified")
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedPaymentStatusCannotBeModified.WithInternal(errStatus),
			method,
			span,
			zap.Error(errStatus),
		)
	}

	_, err = s.merchantQueryRepository.FindById(ctx, int(cashier.MerchantID))
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			merchant_errors.ErrFailedFindMerchantById.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}
	req.MerchantID = int(cashier.MerchantID)

	_, err = s.orderQueryRepository.FindById(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			order_errors.ErrFailedFindOrderById.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			orderitem_errors.ErrFailedFindOrderItemByOrder.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	var totalAmount int
	for _, item := range orderItems {
		totalAmount += int(item.Price) * int(item.Quantity)
	}

	ppn := totalAmount * 11 / 100
	totalAmountWithTax := totalAmount + ppn

	var paymentStatus string
	if req.Amount >= totalAmountWithTax {
		paymentStatus = "success"
	} else {
		status = "error"
		errBalance := errors.New("insufficient balance")
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedPaymentInsufficientBalance.WithInternal(errBalance),
			method,
			span,
			zap.Error(errBalance),
		)
	}

	req.Amount = totalAmountWithTax
	req.PaymentStatus = &paymentStatus

	transaction, err := s.transactionCommandRepository.UpdateTransaction(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedUpdateTransaction.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteTransactionCache(ctx, *req.TransactionID)

	logSuccess("Successfully updated transaction", zap.Bool("success", true))

	return transaction, nil
}

func (s *transactionCommandService) TrashedTransaction(ctx context.Context, transactionID int) (*db.Transaction, error) {
	const method = "TrashedTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("transaction.id", transactionID))

	defer func() {
		end(status)
	}()

	res, err := s.transactionCommandRepository.TrashTransaction(ctx, transactionID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedTrashedTransaction.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteTransactionCache(ctx, transactionID)

	logSuccess("Successfully trashed transaction", zap.Int("transaction.id", transactionID), zap.Bool("success", true))

	return res, nil
}

func (s *transactionCommandService) RestoreTransaction(ctx context.Context, transactionID int) (*db.Transaction, error) {
	const method = "RestoreTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("transaction.id", transactionID))

	defer func() {
		end(status)
	}()

	res, err := s.transactionCommandRepository.RestoreTransaction(ctx, transactionID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedRestoreTransaction.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteTransactionCache(ctx, transactionID)

	logSuccess("Successfully restored transaction", zap.Int("transaction.id", transactionID), zap.Bool("success", true))

	return res, nil
}

func (s *transactionCommandService) DeleteTransactionPermanently(ctx context.Context, transactionID int) (bool, error) {
	const method = "DeleteTransactionPermanently"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("transaction.id", transactionID))

	defer func() {
		end(status)
	}()

	success, err := s.transactionCommandRepository.DeleteTransactionPermanently(ctx, transactionID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			transaction_errors.ErrFailedDeleteTransactionPermanently.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteTransactionCache(ctx, transactionID)

	logSuccess("Successfully permanently deleted transaction", zap.Int("transaction.id", transactionID), zap.Bool("success", success))

	return success, nil
}

func (s *transactionCommandService) RestoreAllTransactions(ctx context.Context) (bool, error) {
	const method = "RestoreAllTransactions"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.transactionCommandRepository.RestoreAllTransactions(ctx)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			transaction_errors.ErrFailedRestoreAllTransactions.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteTransactionAllCache(ctx)
	logSuccess("All trashed transactions restored successfully", zap.Bool("success", success))

	return success, nil
}

func (s *transactionCommandService) DeleteAllTransactionPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllTransactionPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.transactionCommandRepository.DeleteAllTransactionPermanent(ctx)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			transaction_errors.ErrFailedDeleteAllTransactionPermanent.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteTransactionAllCache(ctx)
	logSuccess("Successfully permanently deleted all trashed transactions", zap.Bool("success", success))

	return success, nil
}
