package service

import (
	"context"
	"errors"
	"fmt"

	mencache "github.com/MamangRust/monolith-pointofsale-order/cache"
	"github.com/MamangRust/monolith-pointofsale-order/repository"
	cashieradapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/cashier"
	merchantadapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/merchant"
	orderitemadapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/orderitem"
	productadapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/product"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-pointofsale-shared/errorhandler"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/cashier_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/merchant_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/order_errors"
	orderitem_errors "github.com/MamangRust/monolith-pointofsale-shared/errors/order_item_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/product_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type orderCommandDeps struct {
	Cache                      mencache.OrderCommandCache
	CashierQueryRepository     cashieradapter.QueryRepository
	OrderQueryRepository       repository.OrderQueryRepository
	OrderCommandRepository     repository.OrderCommandRepository
	OrderItemQueryRepository   orderitemadapter.QueryRepository
	OrderItemCommandRepository orderitemadapter.CommandRepository
	MerchantQueryRepository    merchantadapter.Repository
	ProductQueryRepository     productadapter.QueryRepository
	ProductCommandRepository   productadapter.CommandRepository
	Logger                     logger.LoggerInterface
	Observability              observability.TraceLoggerObservability
}

// stockMutation records one applied stock delta so it can be reversed by
// applying the negated delta.
type stockMutation struct {
	productID int
	delta     int
}

type orderCommandService struct {
	mencache                   mencache.OrderCommandCache
	cashierQueryRepository     cashieradapter.QueryRepository
	orderQueryRepository       repository.OrderQueryRepository
	orderCommandRepository     repository.OrderCommandRepository
	orderItemQueryRepository   orderitemadapter.QueryRepository
	orderItemCommandRepository orderitemadapter.CommandRepository
	merchantQueryRepository    merchantadapter.Repository
	productQueryRepository     productadapter.QueryRepository
	productCommandRepository   productadapter.CommandRepository
	logger                     logger.LoggerInterface
	observability              observability.TraceLoggerObservability
}

func NewOrderCommandService(params *orderCommandDeps) OrderCommandService {
	return &orderCommandService{
		mencache:                   params.Cache,
		cashierQueryRepository:     params.CashierQueryRepository,
		orderQueryRepository:       params.OrderQueryRepository,
		orderCommandRepository:     params.OrderCommandRepository,
		orderItemQueryRepository:   params.OrderItemQueryRepository,
		orderItemCommandRepository: params.OrderItemCommandRepository,
		merchantQueryRepository:    params.MerchantQueryRepository,
		productQueryRepository:     params.ProductQueryRepository,
		productCommandRepository:   params.ProductCommandRepository,
		logger:                     params.Logger,
		observability:              params.Observability,
	}
}

func (s *orderCommandService) CreateOrder(ctx context.Context, req *requests.CreateOrderRequest) (*db.Order, error) {
	const method = "CreateOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant.id", req.MerchantID),
		attribute.Int("cashier.id", req.CashierID),
	)
	defer func() {
		end(status)
	}()

	_, err := s.merchantQueryRepository.FindById(ctx, req.MerchantID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](
			s.logger,
			merchant_errors.ErrFailedFindMerchantById.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	_, err = s.cashierQueryRepository.FindById(ctx, req.CashierID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](
			s.logger,
			cashier_errors.ErrFailedFindCashierById.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	products := make([]*db.Product, 0, len(req.Items))
	for _, item := range req.Items {
		product, err := s.productQueryRepository.FindById(ctx, item.ProductID)
		if err != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.Order](s.logger, product_errors.ErrFailedFindProductById.WithInternal(err), method, span, zap.Error(err))
		}
		if item.Quantity <= 0 {
			status = "error"
			return sharederrorhandler.HandleError[*db.Order](s.logger, fmt.Errorf("quantity must be positive"), method, span, zap.Int("product.id", item.ProductID))
		}
		products = append(products, product)
	}

	order, err := s.orderCommandRepository.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		MerchantID: req.MerchantID,
		CashierID:  req.CashierID,
	})
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](s.logger, order_errors.ErrFailedCreateOrder.WithInternal(err), method, span, zap.Error(err))
	}

	span.SetAttributes(attribute.Int("order.id", int(order.OrderID)))
	applied := make([]stockMutation, 0, len(req.Items))
	createdItems := make([]int, 0, len(req.Items))
	compensate := func() {
		for i := len(createdItems) - 1; i >= 0; i-- {
			if rollbackErr := s.orderItemCommandRepository.DeleteOrderItem(ctx, createdItems[i]); rollbackErr != nil {
				s.logger.Error("failed to delete compensated order item", zap.Int("order_item.id", createdItems[i]), zap.Error(rollbackErr))
			}
		}
		if rollbackErr := s.orderCommandRepository.DeleteOrder(ctx, int(order.OrderID)); rollbackErr != nil {
			s.logger.Error("failed to delete compensated order", zap.Int("order.id", int(order.OrderID)), zap.Error(rollbackErr))
		}
		s.reverseStockMutations(ctx, applied)
	}
	for i, item := range req.Items {
		// Reserve through the atomic delta RPC. A stale read cannot cause an
		// oversell; the database guard is the final authority.
		if _, err = s.productCommandRepository.AdjustProductStock(ctx, int(products[i].ProductID), -item.Quantity); err != nil {
			// Preserve the business-level stock error when the atomic guarded
			// UPDATE returned no row because another request consumed the stock.
			if current, lookupErr := s.productQueryRepository.FindById(ctx, item.ProductID); lookupErr == nil && current.CountInStock < int32(item.Quantity) {
				err = order_errors.ErrInsufficientProductStock
			}
			compensate()
			status = "error"
			return sharederrorhandler.HandleError[*db.Order](s.logger, err, method, span, zap.Int("product.id", item.ProductID), zap.Error(err))
		}
		applied = append(applied, stockMutation{productID: int(products[i].ProductID), delta: -item.Quantity})
	}

	for i, item := range req.Items {
		createdItem, createErr := s.orderItemCommandRepository.CreateOrderItem(ctx, &requests.CreateOrderItemRecordRequest{OrderID: int(order.OrderID), ProductID: item.ProductID, Quantity: item.Quantity, Price: int(products[i].Price)})
		if createErr != nil {
			compensate()
			status = "error"
			return sharederrorhandler.HandleError[*db.Order](s.logger, orderitem_errors.ErrFailedCreateOrderItem.WithInternal(createErr), method, span, zap.Error(createErr))
		}
		if createdItem == nil {
			compensate()
			status = "error"
			return sharederrorhandler.HandleError[*db.Order](s.logger, fmt.Errorf("order item repository returned nil item"), method, span)
		}
		createdItems = append(createdItems, int(createdItem.OrderItemID))
	}

	totalPrice, err := s.orderItemQueryRepository.CalculateTotalPrice(ctx, int(order.OrderID))
	if err != nil {
		compensate()
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](
			s.logger,
			orderitem_errors.ErrFailedCalculateTotal.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	res, err := s.orderCommandRepository.UpdateOrder(ctx, &requests.UpdateOrderRecordRequest{
		OrderID:    int(order.OrderID),
		TotalPrice: int(*totalPrice),
	})
	if err != nil {
		compensate()
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](
			s.logger,
			order_errors.ErrFailedUpdateOrder.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	logSuccess("Successfully created order", zap.Int("order.id", int(res.OrderID)))
	return res, nil
}

func (s *orderCommandService) UpdateOrder(ctx context.Context, req *requests.UpdateOrderRequest) (*db.Order, error) {
	const method = "UpdateOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("order.id", *req.OrderID))
	defer func() { end(status) }()

	if _, err := s.orderQueryRepository.FindById(ctx, *req.OrderID); err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](s.logger, order_errors.ErrFailedFindOrderById.WithInternal(err), method, span, zap.Error(err))
	}

	oldItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(ctx, *req.OrderID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](s.logger, orderitem_errors.ErrFailedFindOrderItemByOrder.WithInternal(err), method, span, zap.Error(err))
	}
	oldByID := make(map[int32]*db.OrderItem, len(oldItems))
	for _, item := range oldItems {
		oldByID[item.OrderItemID] = item
	}
	seenItemIDs := make(map[int]struct{}, len(req.Items))

	products := make([]*db.Product, 0, len(req.Items))
	for _, item := range req.Items {
		if item.Quantity <= 0 {
			status = "error"
			return sharederrorhandler.HandleError[*db.Order](s.logger, fmt.Errorf("quantity must be positive"), method, span)
		}
		product, findErr := s.productQueryRepository.FindById(ctx, item.ProductID)
		if findErr != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.Order](s.logger, product_errors.ErrFailedFindProductById.WithInternal(findErr), method, span, zap.Error(findErr))
		}
		products = append(products, product)
	}

	changes := make([]stockMutation, 0, len(req.Items)*2)
	updatedItems := make([]*db.OrderItem, 0, len(req.Items))
	createdItemIDs := make([]int, 0, len(req.Items))
	rollbackEntities := func() {
		for i := len(createdItemIDs) - 1; i >= 0; i-- {
			if rollbackErr := s.orderItemCommandRepository.DeleteOrderItem(ctx, createdItemIDs[i]); rollbackErr != nil {
				s.logger.Error("failed to rollback created order item", zap.Int("order_item.id", createdItemIDs[i]), zap.Error(rollbackErr))
			}
		}
		for i := len(updatedItems) - 1; i >= 0; i-- {
			item := updatedItems[i]
			if _, rollbackErr := s.orderItemCommandRepository.UpdateOrderItem(ctx, &requests.UpdateOrderItemRecordRequest{OrderItemID: int(item.OrderItemID), OrderID: int(item.OrderID), ProductID: int(item.ProductID), Quantity: int(item.Quantity), Price: int(item.Price)}); rollbackErr != nil {
				s.logger.Error("failed to rollback updated order item", zap.Int("order_item.id", int(item.OrderItemID)), zap.Error(rollbackErr))
			}
		}
	}
	rollback := func() {
		rollbackEntities()
		s.reverseStockMutations(ctx, changes)
	}
	apply := func(productID, quantity int, decrement bool) error {
		if quantity <= 0 {
			return nil
		}
		delta := quantity
		if decrement {
			delta = -quantity
		}
		if _, applyErr := s.productCommandRepository.AdjustProductStock(ctx, productID, delta); applyErr != nil {
			if decrement {
				// Preserve the business-level stock error when the atomic guarded
				// UPDATE returned no row because another request consumed the stock.
				if current, lookupErr := s.productQueryRepository.FindById(ctx, productID); lookupErr == nil && current.CountInStock < int32(quantity) {
					return order_errors.ErrInsufficientProductStock
				}
			}
			return applyErr
		}
		changes = append(changes, stockMutation{productID: productID, delta: delta})
		return nil
	}

	for i, item := range req.Items {
		if item.OrderItemID > 0 {
			if _, duplicate := seenItemIDs[item.OrderItemID]; duplicate {
				status = "error"
				return sharederrorhandler.HandleError[*db.Order](s.logger, fmt.Errorf("order item %d appears more than once", item.OrderItemID), method, span)
			}
			seenItemIDs[item.OrderItemID] = struct{}{}
		}
		old := oldByID[int32(item.OrderItemID)]
		if item.OrderItemID > 0 && old == nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.Order](s.logger, fmt.Errorf("order item %d not found", item.OrderItemID), method, span)
		}
		if old != nil {
			if old.ProductID == int32(item.ProductID) {
				delta := item.Quantity - int(old.Quantity)
				if delta > 0 {
					err = apply(item.ProductID, delta, true)
				} else {
					err = apply(item.ProductID, -delta, false)
				}
			} else {
				err = apply(int(old.ProductID), int(old.Quantity), false)
				if err == nil {
					err = apply(item.ProductID, item.Quantity, true)
				}
			}
		} else {
			err = apply(item.ProductID, item.Quantity, true)
		}
		if err != nil {
			rollback()
			status = "error"
			return sharederrorhandler.HandleError[*db.Order](s.logger, err, method, span, zap.Error(err))
		}
		if old != nil {
			if _, err = s.orderItemCommandRepository.UpdateOrderItem(ctx, &requests.UpdateOrderItemRecordRequest{OrderItemID: item.OrderItemID, OrderID: *req.OrderID, ProductID: item.ProductID, Quantity: item.Quantity, Price: int(products[i].Price)}); err != nil {
				rollback()
				status = "error"
				return sharederrorhandler.HandleError[*db.Order](s.logger, orderitem_errors.ErrFailedUpdateOrderItem.WithInternal(err), method, span, zap.Error(err))
			}
			updatedItems = append(updatedItems, old)
		} else {
			createdItem, createErr := s.orderItemCommandRepository.CreateOrderItem(ctx, &requests.CreateOrderItemRecordRequest{OrderID: *req.OrderID, ProductID: item.ProductID, Quantity: item.Quantity, Price: int(products[i].Price)})
			if createErr != nil {
				rollback()
				status = "error"
				return sharederrorhandler.HandleError[*db.Order](s.logger, orderitem_errors.ErrFailedCreateOrderItem.WithInternal(createErr), method, span, zap.Error(createErr))
			}
			if createdItem == nil {
				rollback()
				status = "error"
				return sharederrorhandler.HandleError[*db.Order](s.logger, fmt.Errorf("order item repository returned nil item"), method, span)
			}
			createdItemIDs = append(createdItemIDs, int(createdItem.OrderItemID))
		}
	}

	totalPrice, err := s.orderItemQueryRepository.CalculateTotalPrice(ctx, *req.OrderID)
	if err != nil {
		rollback()
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](s.logger, orderitem_errors.ErrFailedCalculateTotal.WithInternal(err), method, span, zap.Error(err))
	}
	res, err := s.orderCommandRepository.UpdateOrder(ctx, &requests.UpdateOrderRecordRequest{OrderID: *req.OrderID, TotalPrice: int(*totalPrice)})
	if err != nil {
		rollback()
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](s.logger, order_errors.ErrFailedUpdateOrder.WithInternal(err), method, span, zap.Error(err))
	}

	s.mencache.DeleteOrderCache(ctx, *req.OrderID)
	logSuccess("Successfully updated order", zap.Int("order.id", *req.OrderID))
	return res, nil
}

// reverseStockMutations undoes previously applied stock deltas in reverse order
// by applying the negated delta.
func (s *orderCommandService) reverseStockMutations(ctx context.Context, applied []stockMutation) {
	for i := len(applied) - 1; i >= 0; i-- {
		mutation := applied[i]
		if _, err := s.productCommandRepository.AdjustProductStock(ctx, mutation.productID, -mutation.delta); err != nil {
			s.logger.Error("failed to compensate order lifecycle stock", zap.Int("product.id", mutation.productID), zap.Int("delta", mutation.delta), zap.Error(err))
		}
	}
}

func (s *orderCommandService) incrementStockForOrderItems(ctx context.Context, items []*db.OrderItem) ([]stockMutation, error) {
	applied := make([]stockMutation, 0, len(items))
	for _, item := range items {
		if item == nil || item.Quantity <= 0 {
			continue
		}
		if _, err := s.productCommandRepository.AdjustProductStock(ctx, int(item.ProductID), int(item.Quantity)); err != nil {
			s.reverseStockMutations(ctx, applied)
			return nil, err
		}
		applied = append(applied, stockMutation{productID: int(item.ProductID), delta: int(item.Quantity)})
	}
	return applied, nil
}

func (s *orderCommandService) decrementStockForOrderItems(ctx context.Context, items []*db.OrderItem) ([]stockMutation, error) {
	applied := make([]stockMutation, 0, len(items))
	for _, item := range items {
		if item == nil || item.Quantity <= 0 {
			continue
		}
		if _, err := s.productCommandRepository.AdjustProductStock(ctx, int(item.ProductID), -int(item.Quantity)); err != nil {
			s.reverseStockMutations(ctx, applied)
			if current, lookupErr := s.productQueryRepository.FindById(ctx, int(item.ProductID)); lookupErr == nil && current.CountInStock < item.Quantity {
				return nil, order_errors.ErrInsufficientProductStock
			}
			return nil, err
		}
		applied = append(applied, stockMutation{productID: int(item.ProductID), delta: -int(item.Quantity)})
	}
	return applied, nil
}

func (s *orderCommandService) TrashedOrder(ctx context.Context, orderID int) (*db.Order, error) {
	const method = "TrashedOrder"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("order.id", orderID))
	defer func() { end(status) }()

	order, err := s.orderQueryRepository.FindById(ctx, orderID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](s.logger, order_errors.ErrFailedFindOrderById.WithInternal(err), method, span, zap.Error(err))
	}
	if order == nil || order.DeletedAt.Valid {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](s.logger, order_errors.ErrFailedNotDeleteAtOrder, method, span)
	}

	items, err := s.orderItemQueryRepository.FindOrderItemByOrder(ctx, orderID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](s.logger, orderitem_errors.ErrFailedFindOrderItemByOrder.WithInternal(err), method, span, zap.Error(err))
	}
	applied, err := s.incrementStockForOrderItems(ctx, items)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](s.logger, err, method, span, zap.Error(err))
	}

	trashedOrder, err := s.orderCommandRepository.TrashedOrder(ctx, orderID)
	if err != nil || trashedOrder == nil {
		s.reverseStockMutations(ctx, applied)
		status = "error"
		if err == nil {
			err = order_errors.ErrFailedOrderNotActive
		}
		return sharederrorhandler.HandleError[*db.Order](s.logger, order_errors.ErrFailedTrashOrder.WithInternal(err), method, span, zap.Error(err))
	}

	s.mencache.DeleteOrderCache(ctx, orderID)
	logSuccess("Successfully trashed order", zap.Int("order.id", orderID))
	return trashedOrder, nil
}

func (s *orderCommandService) RestoreOrder(ctx context.Context, orderID int) (*db.Order, error) {
	const method = "RestoreOrder"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("order.id", orderID))
	defer func() { end(status) }()

	trashedOrder, err := s.orderQueryRepository.FindByTrashedId(ctx, orderID)
	if err != nil || trashedOrder == nil || !trashedOrder.DeletedAt.Valid {
		status = "error"
		if err == nil {
			err = order_errors.ErrFailedOrderNotTrashed
		}
		return sharederrorhandler.HandleError[*db.Order](s.logger, order_errors.ErrFailedOrderNotTrashed.WithInternal(err), method, span, zap.Error(err))
	}

	items, err := s.orderItemQueryRepository.FindOrderItemByOrder(ctx, orderID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](s.logger, orderitem_errors.ErrFailedFindOrderItemByOrder.WithInternal(err), method, span, zap.Error(err))
	}
	applied, err := s.decrementStockForOrderItems(ctx, items)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](s.logger, err, method, span, zap.Error(err))
	}

	order, err := s.orderCommandRepository.RestoreOrder(ctx, orderID)
	if err != nil || order == nil {
		s.reverseStockMutations(ctx, applied)
		status = "error"
		if err == nil {
			err = order_errors.ErrFailedOrderNotActive
		}
		return sharederrorhandler.HandleError[*db.Order](s.logger, order_errors.ErrFailedRestoreOrder.WithInternal(err), method, span, zap.Error(err))
	}

	s.mencache.DeleteOrderCache(ctx, orderID)
	logSuccess("Successfully restored order", zap.Int("order.id", orderID))
	return order, nil
}

func (s *orderCommandService) DeleteOrderPermanent(ctx context.Context, orderID int) (bool, error) {
	const method = "DeleteOrderPermanent"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("order.id", orderID))
	defer func() { end(status) }()

	order, err := s.orderQueryRepository.FindByTrashedId(ctx, orderID)
	if err != nil || order == nil || !order.DeletedAt.Valid {
		status = "error"
		if err == nil {
			err = order_errors.ErrFailedOrderNotTrashed
		}
		return sharederrorhandler.HandleError[bool](s.logger, order_errors.ErrFailedOrderNotTrashed.WithInternal(err), method, span, zap.Error(err))
	}

	success, err := s.orderCommandRepository.DeleteOrderPermanent(ctx, orderID)
	if err != nil || !success {
		status = "error"
		if err == nil {
			err = order_errors.ErrFailedDeleteOrderPermanent
		}
		return sharederrorhandler.HandleError[bool](s.logger, order_errors.ErrFailedDeleteOrderPermanent.WithInternal(err), method, span, zap.Error(err))
	}

	s.mencache.DeleteOrderCache(ctx, orderID)
	logSuccess("Successfully deleted order permanently", zap.Int("order.id", orderID))
	return success, nil
}

func (s *orderCommandService) RestoreAllOrder(ctx context.Context) (bool, error) {
	const method = "RestoreAllOrder"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() { end(status) }()

	orders, err := s.orderCommandRepository.FindAllTrashed(ctx)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](s.logger, order_errors.ErrFailedRestoreAllOrder.WithInternal(err), method, span, zap.Error(err))
	}

	type restoredOrder struct {
		orderID int
		items   []*db.OrderItem
		stock   []stockMutation
	}
	restoredOrders := make([]restoredOrder, 0, len(orders))
	rollbackAll := func() {
		for i := len(restoredOrders) - 1; i >= 0; i-- {
			candidate := restoredOrders[i]
			if _, rollbackErr := s.orderCommandRepository.TrashedOrder(ctx, candidate.orderID); rollbackErr != nil {
				s.logger.Error("failed to rollback restored order state", zap.Int("order.id", candidate.orderID), zap.Error(rollbackErr))
			}
			s.reverseStockMutations(ctx, candidate.stock)
		}
	}

	restoredCount := 0
	for _, candidate := range orders {
		if candidate == nil || !candidate.DeletedAt.Valid {
			continue
		}
		items, findErr := s.orderItemQueryRepository.FindOrderItemByOrder(ctx, int(candidate.OrderID))
		if findErr != nil {
			rollbackAll()
			status = "error"
			return sharederrorhandler.HandleError[bool](s.logger, order_errors.ErrFailedRestoreAllOrder.WithInternal(findErr), method, span, zap.Error(findErr))
		}
		applied, decrementErr := s.decrementStockForOrderItems(ctx, items)
		if decrementErr != nil {
			rollbackAll()
			status = "error"
			return sharederrorhandler.HandleError[bool](s.logger, order_errors.ErrFailedRestoreAllOrder.WithInternal(decrementErr), method, span, zap.Error(decrementErr))
		}
		restored, restoreErr := s.orderCommandRepository.RestoreOrder(ctx, int(candidate.OrderID))
		if restoreErr != nil || restored == nil {
			s.reverseStockMutations(ctx, applied)
			if restoreErr != nil && errors.Is(restoreErr, order_errors.ErrRestoreOrderNotFound) {
				continue
			}
			rollbackAll()
			if restoreErr == nil {
				restoreErr = order_errors.ErrFailedOrderNotActive
			}
			status = "error"
			return sharederrorhandler.HandleError[bool](s.logger, order_errors.ErrFailedRestoreAllOrder.WithInternal(restoreErr), method, span, zap.Error(restoreErr))
		}
		restoredOrders = append(restoredOrders, restoredOrder{orderID: int(candidate.OrderID), items: items, stock: applied})
		restoredCount++
	}

	s.mencache.DeleteOrderAllCache(ctx)
	logSuccess("Successfully restored all orders", zap.Int("restored.count", restoredCount))
	return true, nil
}

func (s *orderCommandService) DeleteAllOrderPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllOrderPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	success, err := s.orderCommandRepository.DeleteAllOrderPermanent(ctx)
	if err != nil || !success {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			order_errors.ErrFailedDeleteAllOrderPermanent.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteOrderAllCache(ctx)
	logSuccess("Successfully deleted all orders permanently", zap.Bool("success", success))
	return success, nil
}
