package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-pointofsale-cashier/cache"
	"github.com/MamangRust/monolith-pointofsale-cashier/repository"
	merchantadapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/merchant"
	useradapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/user"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-pointofsale-shared/errorhandler"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/cashier_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/merchant_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/user_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type cashierCommandDeps struct {
	Cache          mencache.CashierCommandCache
	MerchantQuery  merchantadapter.Repository
	UserQuery      useradapter.QueryRepository
	CashierCommand repository.CashierCommandRepository
	Logger         logger.LoggerInterface
	Observability  observability.TraceLoggerObservability
}

type cashierCommandService struct {
	mencache       mencache.CashierCommandCache
	merchantQuery  merchantadapter.Repository
	userQuery      useradapter.QueryRepository
	cashierCommand repository.CashierCommandRepository
	logger         logger.LoggerInterface
	observability  observability.TraceLoggerObservability
}

func NewCashierCommandService(params *cashierCommandDeps) CashierCommandService {
	return &cashierCommandService{
		mencache:       params.Cache,
		merchantQuery:  params.MerchantQuery,
		userQuery:      params.UserQuery,
		cashierCommand: params.CashierCommand,
		logger:         params.Logger,
		observability:  params.Observability,
	}
}

func (s *cashierCommandService) CreateCashier(ctx context.Context, req *requests.CreateCashierRequest) (*db.Cashier, error) {
	const method = "CreateCashier"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	_, err := s.merchantQuery.FindById(ctx, req.MerchantID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Cashier](
			s.logger,
			merchant_errors.ErrFailedFindMerchantById,
			method,
			span,
			zap.Error(err),
		)
	}

	_, err = s.userQuery.FindById(ctx, req.UserID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Cashier](
			s.logger,
			user_errors.ErrUserNotFoundRes,
			method,
			span,
			zap.Error(err),
		)
	}

	cashier, err := s.cashierCommand.CreateCashier(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Cashier](
			s.logger,
			cashier_errors.ErrFailedCreateCashier,
			method,
			span,
			zap.Error(err),
		)
	}

	logSuccess("Successfully created cashier", zap.Bool("success", true))
	return cashier, nil
}

func (s *cashierCommandService) UpdateCashier(ctx context.Context, req *requests.UpdateCashierRequest) (*db.Cashier, error) {
	const method = "UpdateCashier"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	cashier, err := s.cashierCommand.UpdateCashier(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Cashier](
			s.logger,
			cashier_errors.ErrFailedUpdateCashier,
			method,
			span,
			zap.Error(err),
		)
	}

	span.SetAttributes(attribute.String("cashier.name", cashier.Name))
	s.mencache.DeleteCashierCache(ctx, int(cashier.CashierID))

	logSuccess("Successfully updated cashier", zap.Bool("success", true))
	return cashier, nil
}

func (s *cashierCommandService) TrashedCashier(ctx context.Context, cashierID int) (*db.Cashier, error) {
	const method = "TrashedCashier"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	cashier, err := s.cashierCommand.TrashedCashier(ctx, cashierID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Cashier](
			s.logger,
			cashier_errors.ErrFailedTrashedCashier,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteCashierCache(ctx, cashierID)
	logSuccess("Successfully trashed cashier", zap.Bool("success", true))
	return cashier, nil
}

func (s *cashierCommandService) RestoreCashier(ctx context.Context, cashierID int) (*db.Cashier, error) {
	const method = "RestoreCashier"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	cashier, err := s.cashierCommand.RestoreCashier(ctx, cashierID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Cashier](
			s.logger,
			cashier_errors.ErrFailedRestoreCashier,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteCashierCache(ctx, cashierID)
	logSuccess("Successfully restored cashier", zap.Bool("success", true))
	return cashier, nil
}

func (s *cashierCommandService) DeleteCashierPermanent(ctx context.Context, cashierID int) (bool, error) {
	const method = "DeleteCashierPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	success, err := s.cashierCommand.DeleteCashierPermanent(ctx, cashierID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			cashier_errors.ErrFailedDeleteCashierPermanent,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteCashierCache(ctx, cashierID)
	logSuccess("Successfully deleted cashier permanently", zap.Bool("success", success))
	return success, nil
}

func (s *cashierCommandService) RestoreAllCashier(ctx context.Context) (bool, error) {
	const method = "RestoreAllCashier"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	success, err := s.cashierCommand.RestoreAllCashier(ctx)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			cashier_errors.ErrFailedRestoreAllCashiers,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteCashierListCache(ctx)
	logSuccess("Successfully restored all trashed cashiers", zap.Bool("success", success))
	return success, nil
}

func (s *cashierCommandService) DeleteAllCashierPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllCashierPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	success, err := s.cashierCommand.DeleteAllCashierPermanent(ctx)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			cashier_errors.ErrFailedDeleteAllCashierPermanent,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteCashierListCache(ctx)
	logSuccess("Successfully deleted all trashed cashiers", zap.Bool("success", success))
	return success, nil
}
