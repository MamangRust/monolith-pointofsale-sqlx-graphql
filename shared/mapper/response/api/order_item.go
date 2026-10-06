package response_api

import (
	pb "github.com/MamangRust/monolith-pointofsale-pb"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"
)

type orderItemResponseMapper struct {
}

func NewOrderItemResponseMapper() *orderItemResponseMapper {
	return &orderItemResponseMapper{}
}

func (o *orderItemResponseMapper) ToResponseOrderItem(orderItem *pb.OrderItemResponse) *response.OrderItemResponse {
	return &response.OrderItemResponse{
		ID:        int(orderItem.Id),
		OrderID:   int(orderItem.OrderId),
		ProductID: int(orderItem.ProductId),
		Quantity:  int(orderItem.Quantity),
		Price:     int(orderItem.Price),
		CreatedAt: orderItem.CreatedAt,
		UpdatedAt: orderItem.UpdatedAt,
	}
}

func (o *orderItemResponseMapper) ToResponsesOrderItem(orderItems []*pb.OrderItemResponse) []*response.OrderItemResponse {
	var mappedOrderItems []*response.OrderItemResponse

	for _, orderItem := range orderItems {
		mappedOrderItems = append(mappedOrderItems, o.ToResponseOrderItem(orderItem))
	}

	return mappedOrderItems
}

func (o *orderItemResponseMapper) ToResponseOrderItemDeleteAt(orderItem *pb.OrderItemResponseDeleteAt) *response.OrderItemResponseDeleteAt {
	var deletedAt string
	if orderItem.DeletedAt != nil {
		deletedAt = orderItem.DeletedAt.Value
	}
	return &response.OrderItemResponseDeleteAt{
		ID:        int(orderItem.Id),
		OrderID:   int(orderItem.OrderId),
		ProductID: int(orderItem.ProductId),
		Quantity:  int(orderItem.Quantity),
		Price:     int(orderItem.Price),
		CreatedAt: orderItem.CreatedAt,
		UpdatedAt: orderItem.UpdatedAt,
		DeleteAt:  &deletedAt,
	}
}

func (o *orderItemResponseMapper) ToResponsesOrderItemDeleteAt(orderItems []*pb.OrderItemResponseDeleteAt) []*response.OrderItemResponseDeleteAt {
	var mappedOrderItems []*response.OrderItemResponseDeleteAt

	for _, orderItem := range orderItems {
		mappedOrderItems = append(mappedOrderItems, o.ToResponseOrderItemDeleteAt(orderItem))
	}

	return mappedOrderItems
}

func (o *orderItemResponseMapper) ToApiResponseOrderItem(pbResponse *pb.ApiResponseOrderItem) *response.ApiResponseOrderItem {
	return &response.ApiResponseOrderItem{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponseOrderItem(pbResponse.Data),
	}
}

func (o *orderItemResponseMapper) ToApiResponsesOrderItem(pbResponse *pb.ApiResponsesOrderItem) *response.ApiResponsesOrderItem {
	return &response.ApiResponsesOrderItem{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponsesOrderItem(pbResponse.Data),
	}
}

func (o *orderItemResponseMapper) ToApiResponseOrderItemDelete(pbResponse *pb.ApiResponseOrderItemDelete) *response.ApiResponseOrderItemDelete {
	return &response.ApiResponseOrderItemDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (o *orderItemResponseMapper) ToApiResponseOrderItemAll(pbResponse *pb.ApiResponseOrderItemAll) *response.ApiResponseOrderItemAll {
	return &response.ApiResponseOrderItemAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (o *orderItemResponseMapper) ToApiResponsePaginationOrderItemDeleteAt(pbResponse *pb.ApiResponsePaginationOrderItemDeleteAt) *response.ApiResponsePaginationOrderItemDeleteAt {
	return &response.ApiResponsePaginationOrderItemDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       o.ToResponsesOrderItemDeleteAt(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (o *orderItemResponseMapper) ToApiResponsePaginationOrderItem(pbResponse *pb.ApiResponsePaginationOrderItem) *response.ApiResponsePaginationOrderItem {
	return &response.ApiResponsePaginationOrderItem{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       o.ToResponsesOrderItem(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}
