package orderitem_errors

import (
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
)

var (
	ErrFindAllOrderItems        = errors.ErrInternal.WithMessage("Failed to find all order items")
	ErrFindByActive             = errors.ErrInternal.WithMessage("Failed to find active order items")
	ErrFindByTrashed            = errors.ErrInternal.WithMessage("Failed to find trashed order items")
	ErrFindOrderItemByOrder     = errors.ErrInternal.WithMessage("Failed to find order items by order ID")
	ErrCalculateTotalPrice      = errors.ErrInternal.WithMessage("Failed to calculate total price")
	ErrCreateOrderItem          = errors.ErrInternal.WithMessage("Failed to create order item")
	ErrUpdateOrderItem          = errors.ErrInternal.WithMessage("Failed to update order item")
	ErrTrashedOrderItem         = errors.ErrInternal.WithMessage("Failed to move order item to trash")
	ErrRestoreOrderItem         = errors.ErrInternal.WithMessage("Failed to restore order item from trash")
	ErrDeleteOrderItemPermanent = errors.ErrInternal.WithMessage("Failed to permanently delete order item")
	ErrDeleteOrderItemsByOrder  = errors.ErrInternal.WithMessage("Failed to permanently delete order items by order")
	ErrRestoreAllOrderItem      = errors.ErrInternal.WithMessage("Failed to restore all trashed order items")
	ErrDeleteAllOrderPermanent  = errors.ErrInternal.WithMessage("Failed to permanently delete all trashed order items")
)
