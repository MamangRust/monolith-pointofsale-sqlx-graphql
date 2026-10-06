package handler

import (
	"context"
	pb "github.com/MamangRust/monolith-pointofsale-pb"
	"math"

	"github.com/MamangRust/monolith-pointofsale-order-item/service"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
	orderitem_errors "github.com/MamangRust/monolith-pointofsale-shared/errors/order_item_errors"

	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type orderItemHandleGrpc struct {
	pb.UnimplementedOrderItemQueryServiceServer
	orderItemService service.OrderItemQueryService
	logger           logger.LoggerInterface
}

func NewOrderItemHandlerGrpc(
	orderItemService service.OrderItemQueryService,
	logger logger.LoggerInterface,
) OrderItemHandlerGrpc {
	return &orderItemHandleGrpc{
		orderItemService: orderItemService,
		logger:           logger,
	}
}

func (s *orderItemHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllOrderItemRequest) (*pb.ApiResponsePaginationOrderItem, error) {
	s.logger.Info("FindAll order items called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrderItems{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orderItems, totalRecords, err := s.orderItemService.FindAllOrderItems(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll order items failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAll order items success")

	return &pb.ApiResponsePaginationOrderItem{
		Status:     "success",
		Message:    "Successfully fetched order items",
		Data:       mapResponsesOrderItem(orderItems),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *orderItemHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllOrderItemRequest) (*pb.ApiResponsePaginationOrderItemDeleteAt, error) {
	s.logger.Info("FindByActive order items called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrderItems{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orderItems, totalRecords, err := s.orderItemService.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive order items failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByActive order items success")

	return &pb.ApiResponsePaginationOrderItemDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active order items",
		Data:       mapResponsesOrderItemActive(orderItems),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *orderItemHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllOrderItemRequest) (*pb.ApiResponsePaginationOrderItemDeleteAt, error) {
	s.logger.Info("FindByTrashed order items called")

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrderItems{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orderItems, totalRecords, err := s.orderItemService.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed order items failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByTrashed order items success")

	return &pb.ApiResponsePaginationOrderItemDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed order items",
		Data:       mapResponsesOrderItemTrashed(orderItems),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *orderItemHandleGrpc) FindOrderItemByOrder(ctx context.Context, request *pb.FindByIdOrderItemRequest) (*pb.ApiResponsesOrderItem, error) {
	s.logger.Info("FindOrderItemByOrder called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, orderitem_errors.ErrGrpcInvalidID
	}

	orderItems, err := s.orderItemService.FindOrderItemByOrder(ctx, id)
	if err != nil {
		s.logger.Error("FindOrderItemByOrder failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindOrderItemByOrder success")

	return &pb.ApiResponsesOrderItem{
		Status:  "success",
		Message: "Successfully fetched order items by order",
		Data:    mapResponsesOrderItemFromModel(orderItems),
	}, nil
}

// Map helpers
func mapPaginationMeta(meta *pbcommon.PaginationMeta) *pbcommon.PaginationMeta {
	if meta == nil {
		return nil
	}
	return &pbcommon.PaginationMeta{
		CurrentPage:  meta.CurrentPage,
		PageSize:     meta.PageSize,
		TotalPages:   meta.TotalPages,
		TotalRecords: meta.TotalRecords,
	}
}

func mapGetOrderItemsRowToProto(orderItem *db.GetOrderItemsRow) *pb.OrderItemResponse {
	if orderItem == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	if orderItem.CreatedAt.Valid {
		createdAtStr = orderItem.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if orderItem.UpdatedAt.Valid {
		updatedAtStr = orderItem.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	return &pb.OrderItemResponse{
		Id:        orderItem.OrderItemID,
		OrderId:   orderItem.OrderID,
		ProductId: orderItem.ProductID,
		Quantity:  orderItem.Quantity,
		Price:     orderItem.Price,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
	}
}

func mapResponsesOrderItem(orderItems []*db.GetOrderItemsRow) []*pb.OrderItemResponse {
	var mappedOrderItems []*pb.OrderItemResponse
	for _, orderItem := range orderItems {
		mappedOrderItems = append(mappedOrderItems, mapGetOrderItemsRowToProto(orderItem))
	}
	return mappedOrderItems
}

func mapGetOrderItemsActiveRowToProto(orderItem *db.GetOrderItemsActiveRow) *pb.OrderItemResponseDeleteAt {
	if orderItem == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	var deletedAt *wrapperspb.StringValue
	if orderItem.CreatedAt.Valid {
		createdAtStr = orderItem.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if orderItem.UpdatedAt.Valid {
		updatedAtStr = orderItem.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if orderItem.DeletedAt.Valid {
		deletedAt = wrapperspb.String(orderItem.DeletedAt.Time.Format("2006-01-02 15:04:05"))
	}

	return &pb.OrderItemResponseDeleteAt{
		Id:        orderItem.OrderItemID,
		OrderId:   orderItem.OrderID,
		ProductId: orderItem.ProductID,
		Quantity:  orderItem.Quantity,
		Price:     orderItem.Price,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
		DeletedAt: deletedAt,
	}
}

func mapResponsesOrderItemActive(orderItems []*db.GetOrderItemsActiveRow) []*pb.OrderItemResponseDeleteAt {
	var mappedOrderItems []*pb.OrderItemResponseDeleteAt
	for _, orderItem := range orderItems {
		mappedOrderItems = append(mappedOrderItems, mapGetOrderItemsActiveRowToProto(orderItem))
	}
	return mappedOrderItems
}

func mapGetOrderItemsTrashedRowToProto(orderItem *db.GetOrderItemsTrashedRow) *pb.OrderItemResponseDeleteAt {
	if orderItem == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	var deletedAt *wrapperspb.StringValue
	if orderItem.CreatedAt.Valid {
		createdAtStr = orderItem.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if orderItem.UpdatedAt.Valid {
		updatedAtStr = orderItem.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if orderItem.DeletedAt.Valid {
		deletedAt = wrapperspb.String(orderItem.DeletedAt.Time.Format("2006-01-02 15:04:05"))
	}

	return &pb.OrderItemResponseDeleteAt{
		Id:        orderItem.OrderItemID,
		OrderId:   orderItem.OrderID,
		ProductId: orderItem.ProductID,
		Quantity:  orderItem.Quantity,
		Price:     orderItem.Price,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
		DeletedAt: deletedAt,
	}
}

func mapResponsesOrderItemTrashed(orderItems []*db.GetOrderItemsTrashedRow) []*pb.OrderItemResponseDeleteAt {
	var mappedOrderItems []*pb.OrderItemResponseDeleteAt
	for _, orderItem := range orderItems {
		mappedOrderItems = append(mappedOrderItems, mapGetOrderItemsTrashedRowToProto(orderItem))
	}
	return mappedOrderItems
}

func mapOrderItemToProto(orderItem *db.OrderItem) *pb.OrderItemResponse {
	if orderItem == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	if orderItem.CreatedAt.Valid {
		createdAtStr = orderItem.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if orderItem.UpdatedAt.Valid {
		updatedAtStr = orderItem.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	return &pb.OrderItemResponse{
		Id:        orderItem.OrderItemID,
		OrderId:   orderItem.OrderID,
		ProductId: orderItem.ProductID,
		Quantity:  orderItem.Quantity,
		Price:     orderItem.Price,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
	}
}

func mapResponsesOrderItemFromModel(orderItems []*db.OrderItem) []*pb.OrderItemResponse {
	var mappedOrderItems []*pb.OrderItemResponse
	for _, orderItem := range orderItems {
		mappedOrderItems = append(mappedOrderItems, mapOrderItemToProto(orderItem))
	}
	return mappedOrderItems
}
