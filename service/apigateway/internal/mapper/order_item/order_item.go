package orderitemgraphqlmapper

import (
	graphqlmapper "github.com/MamangRust/monolith-pointofsale-apigateway/internal/mapper"
	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
	pb "github.com/MamangRust/monolith-pointofsale-pb"
)

type orderItemGraphqlMapper struct {
}

func NewOrderItemGraphqlMapper() *orderItemGraphqlMapper {
	return &orderItemGraphqlMapper{}
}

func (o *orderItemGraphqlMapper) ToGraphqlResponseOrderItem(res *pb.ApiResponseOrderItem) *model.APIResponseOrderItem {
	return &model.APIResponseOrderItem{
		Status:  res.Status,
		Message: res.Message,
		Data:    o.mapResponseOrderItem(res.Data),
	}
}

func (o *orderItemGraphqlMapper) ToGraphqlResponsesOrderItem(res *pb.ApiResponsesOrderItem) *model.APIResponsesOrderItem {
	return &model.APIResponsesOrderItem{
		Status:  res.Status,
		Message: res.Message,
		Data:    o.mapResponsesOrderItem(res.Data),
	}
}

func (o *orderItemGraphqlMapper) ToGrapqhlResponseOrderItemDelete(res *pb.ApiResponseOrderItemDelete) *model.APIResponseOrderItemDelete {
	return &model.APIResponseOrderItemDelete{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (o *orderItemGraphqlMapper) ToGrapqhlResponseOrderItemAll(res *pb.ApiResponseOrderItemAll) *model.APIResponseOrderItemAll {
	return &model.APIResponseOrderItemAll{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (o *orderItemGraphqlMapper) ToGraphqlResponsePaginationOrderItem(res *pb.ApiResponsePaginationOrderItem) *model.APIResponsePaginationOrderItem {
	return &model.APIResponsePaginationOrderItem{
		Status:     res.Status,
		Message:    res.Message,
		Data:       o.mapResponsesOrderItem(res.Data),
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}

func (o *orderItemGraphqlMapper) ToGraphqlResponsePaginationOrderItemDeleteAt(res *pb.ApiResponsePaginationOrderItemDeleteAt) *model.APIResponsePaginationOrderItemDeleteAt {
	return &model.APIResponsePaginationOrderItemDeleteAt{
		Status:     res.Status,
		Message:    res.Message,
		Data:       o.mapResponsesOrderItemDeleteAt(res.Data),
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}

func (o *orderItemGraphqlMapper) mapResponseOrderItem(orderItem *pb.OrderItemResponse) *model.OrderItemResponse {
	if orderItem == nil {
		return nil
	}
	return &model.OrderItemResponse{
		ID:        int32(orderItem.Id),
		OrderID:   int32(orderItem.OrderId),
		ProductID: int32(orderItem.ProductId),
		Quantity:  int32(orderItem.Quantity),
		Price:     int32(orderItem.Price),
		CreatedAt: &orderItem.CreatedAt,
		UpdatedAt: &orderItem.UpdatedAt,
	}
}

func (o *orderItemGraphqlMapper) mapResponsesOrderItem(orderItems []*pb.OrderItemResponse) []*model.OrderItemResponse {
	var responses []*model.OrderItemResponse
	for _, orderitem := range orderItems {
		responses = append(responses, o.mapResponseOrderItem(orderitem))
	}
	return responses
}

func (o *orderItemGraphqlMapper) mapResponseOrderItemDeleteAt(orderItem *pb.OrderItemResponseDeleteAt) *model.OrderItemResponseDeleteAt {
	if orderItem == nil {
		return nil
	}
	var deletedAt string
	if orderItem.DeletedAt != nil {
		deletedAt = orderItem.DeletedAt.Value
	}

	return &model.OrderItemResponseDeleteAt{
		ID:        int32(orderItem.Id),
		OrderID:   int32(orderItem.OrderId),
		ProductID: int32(orderItem.ProductId),
		Quantity:  int32(orderItem.Quantity),
		Price:     int32(orderItem.Price),
		CreatedAt: &orderItem.CreatedAt,
		UpdatedAt: &orderItem.UpdatedAt,
		DeletedAt: &deletedAt,
	}
}

func (o *orderItemGraphqlMapper) mapResponsesOrderItemDeleteAt(orderItems []*pb.OrderItemResponseDeleteAt) []*model.OrderItemResponseDeleteAt {
	var responses []*model.OrderItemResponseDeleteAt
	for _, orderitem := range orderItems {
		responses = append(responses, o.mapResponseOrderItemDeleteAt(orderitem))
	}
	return responses
}
