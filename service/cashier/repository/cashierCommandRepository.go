package repository

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/cashier_errors"
	"github.com/jmoiron/sqlx"
)

const createCashier = `INSERT INTO cashiers (merchant_id, user_id, name)
VALUES ($1, $2, $3) RETURNING *`

const updateCashier = `UPDATE cashiers
SET name = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE cashier_id = $1
  AND deleted_at IS NULL
  RETURNING *`

const trashCashier = `UPDATE cashiers
SET
    deleted_at = current_timestamp
WHERE
    cashier_id = $1
    AND deleted_at IS NULL
    RETURNING *`

const restoreCashier = `UPDATE cashiers
SET
    deleted_at = NULL
WHERE
    cashier_id = $1
    AND deleted_at IS NOT NULL
  RETURNING *`

const deleteCashierPermanently = `DELETE FROM cashiers WHERE cashier_id = $1 AND deleted_at IS NOT NULL`

const restoreAllCashiers = `UPDATE cashiers
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL`

const deleteAllPermanentCashiers = `DELETE FROM cashiers
WHERE
    deleted_at IS NOT NULL`

type cashierCommandRepository struct {
	db *sqlx.DB
}

func NewCashierCommandRepository(db *sqlx.DB) CashierCommandRepository {
	return &cashierCommandRepository{
		db: db,
	}
}

func (r *cashierCommandRepository) CreateCashier(ctx context.Context, request *requests.CreateCashierRequest) (*db.Cashier, error) {
	var cashier db.Cashier
	err := r.db.GetContext(ctx, &cashier, createCashier,
		int32(request.MerchantID),
		int32(request.UserID),
		request.Name,
	)
	if err != nil {
		return nil, cashier_errors.ErrCreateCashier
	}

	return &cashier, nil
}

func (r *cashierCommandRepository) UpdateCashier(ctx context.Context, request *requests.UpdateCashierRequest) (*db.Cashier, error) {
	var cashier db.Cashier
	err := r.db.GetContext(ctx, &cashier, updateCashier,
		int32(*request.CashierID),
		request.Name,
	)
	if err != nil {
		return nil, cashier_errors.ErrUpdateCashier
	}

	return &cashier, nil
}

func (r *cashierCommandRepository) TrashedCashier(ctx context.Context, cashier_id int) (*db.Cashier, error) {
	var cashier db.Cashier
	err := r.db.GetContext(ctx, &cashier, trashCashier, int32(cashier_id))
	if err != nil {
		return nil, cashier_errors.ErrTrashedCashier
	}

	return &cashier, nil
}

func (r *cashierCommandRepository) RestoreCashier(ctx context.Context, cashier_id int) (*db.Cashier, error) {
	var cashier db.Cashier
	err := r.db.GetContext(ctx, &cashier, restoreCashier, int32(cashier_id))
	if err != nil {
		return nil, cashier_errors.ErrRestoreCashier
	}

	return &cashier, nil
}

func (r *cashierCommandRepository) DeleteCashierPermanent(ctx context.Context, cashier_id int) (bool, error) {
	if _, err := r.db.ExecContext(ctx, deleteCashierPermanently, int32(cashier_id)); err != nil {
		return false, cashier_errors.ErrDeleteCashierPermanent
	}

	return true, nil
}

func (r *cashierCommandRepository) RestoreAllCashier(ctx context.Context) (bool, error) {
	if _, err := r.db.ExecContext(ctx, restoreAllCashiers); err != nil {
		return false, cashier_errors.ErrRestoreAllCashiers
	}

	return true, nil
}

func (r *cashierCommandRepository) DeleteAllCashierPermanent(ctx context.Context) (bool, error) {
	if _, err := r.db.ExecContext(ctx, deleteAllPermanentCashiers); err != nil {
		return false, cashier_errors.ErrDeleteAllCashiersPermanent
	}

	return true, nil
}
