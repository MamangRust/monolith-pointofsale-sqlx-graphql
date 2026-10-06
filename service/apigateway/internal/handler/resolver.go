package graph

import (
	errorstd "errors"
	"fmt"
	"time"

	authgraphqlmapper "github.com/MamangRust/monolith-pointofsale-apigateway/internal/mapper/auth"
	cashiergraphqlmapper "github.com/MamangRust/monolith-pointofsale-apigateway/internal/mapper/cashier"
	categorygraphqlmapper "github.com/MamangRust/monolith-pointofsale-apigateway/internal/mapper/category"
	merchantgraphqlmapper "github.com/MamangRust/monolith-pointofsale-apigateway/internal/mapper/merchant"
	merchantdocumentgraphqlmapper "github.com/MamangRust/monolith-pointofsale-apigateway/internal/mapper/merchant_document"
	ordergraphqlmapper "github.com/MamangRust/monolith-pointofsale-apigateway/internal/mapper/order"
	orderitemgraphqlmapper "github.com/MamangRust/monolith-pointofsale-apigateway/internal/mapper/order_item"
	productgraphqlmapper "github.com/MamangRust/monolith-pointofsale-apigateway/internal/mapper/product"
	rolegraphqlmapper "github.com/MamangRust/monolith-pointofsale-apigateway/internal/mapper/role"
	transactiongraphqlmapper "github.com/MamangRust/monolith-pointofsale-apigateway/internal/mapper/transaction"
	usergraphqlmapper "github.com/MamangRust/monolith-pointofsale-apigateway/internal/mapper/user"

	merchantpermission "github.com/MamangRust/monolith-pointofsale-apigateway/internal/permission/merchant"
	rolepermission "github.com/MamangRust/monolith-pointofsale-apigateway/internal/permission/role"
	mencache "github.com/MamangRust/monolith-pointofsale-apigateway/internal/redis"

	auth_cache "github.com/MamangRust/monolith-pointofsale-apigateway/internal/redis/api/auth"
	cashier_cache "github.com/MamangRust/monolith-pointofsale-apigateway/internal/redis/api/cashier"
	category_cache "github.com/MamangRust/monolith-pointofsale-apigateway/internal/redis/api/category"
	merchant_cache "github.com/MamangRust/monolith-pointofsale-apigateway/internal/redis/api/merchant"
	merchant_document_cache "github.com/MamangRust/monolith-pointofsale-apigateway/internal/redis/api/merchant_document"
	order_cache "github.com/MamangRust/monolith-pointofsale-apigateway/internal/redis/api/order"
	orderitem_cache "github.com/MamangRust/monolith-pointofsale-apigateway/internal/redis/api/order_item"
	product_cache "github.com/MamangRust/monolith-pointofsale-apigateway/internal/redis/api/product"
	role_cache "github.com/MamangRust/monolith-pointofsale-apigateway/internal/redis/api/role"
	transaction_cache "github.com/MamangRust/monolith-pointofsale-apigateway/internal/redis/api/transaction"
	user_cache "github.com/MamangRust/monolith-pointofsale-apigateway/internal/redis/api/user"

	authpb "github.com/MamangRust/monolith-pointofsale-pb"
	orderitempb "github.com/MamangRust/monolith-pointofsale-pb"
	cashierpb "github.com/MamangRust/monolith-pointofsale-pb/cashier"

	categorypb "github.com/MamangRust/monolith-pointofsale-pb/category"
	merchantpb "github.com/MamangRust/monolith-pointofsale-pb/merchant"
	merchantdocumentpb "github.com/MamangRust/monolith-pointofsale-pb/merchant_document"
	orderpb "github.com/MamangRust/monolith-pointofsale-pb/order"
	productpb "github.com/MamangRust/monolith-pointofsale-pb/product"
	rolepb "github.com/MamangRust/monolith-pointofsale-pb/role"
	transactionpb "github.com/MamangRust/monolith-pointofsale-pb/transaction"
	userpb "github.com/MamangRust/monolith-pointofsale-pb/user"

	"github.com/MamangRust/monolith-pointofsale-pkg/kafka"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-pkg/upload_image"
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
	sharedErrors "github.com/MamangRust/monolith-pointofsale-shared/errors"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AuthGraphql             AuthHandleGraphql
	RoleGraphql             RoleHandleGraphql
	UserGraphql             UserHandleGraphql
	CashierGraphql          CashierHandleGraphql
	CategoryGraphql         CategoryHandleGraphql
	MerchantGraphql         MerchantHandleGraphql
	MerchantDocumentGraphql MerchantDocumentHandleGraphql
	OrderGraphql            OrderHandleGraphql
	OrderItemGraphql        OrderItemHandleGraphql
	ProductGraphql          ProductHandleGraphql
	TransactionGraphql      TransactionHandleGraphql
	ResolverHandle          *resolverHandler
}

type UserClient struct {
	UserQueryClient   userpb.UserQueryServiceClient
	UserCommandClient userpb.UserCommandServiceClient
}

type RoleClient struct {
	rolepb.RoleQueryServiceClient
	rolepb.RoleCommandServiceClient
}

type CashierClient struct {
	cashierpb.CashierQueryServiceClient
	cashierpb.CashierCommandServiceClient
	cashierpb.CashierStatsServiceClient
	cashierpb.CashierStatsByIdServiceClient
	cashierpb.CashierStatsByMerchantServiceClient
}

type CategoryClient struct {
	categorypb.CategoryQueryServiceClient
	categorypb.CategoryCommandServiceClient
	categorypb.CategoryStatsServiceClient
	categorypb.CategoryStatsByIdServiceClient
	categorypb.CategoryStatsByMerchantServiceClient
}

type MerchantClient struct {
	merchantpb.MerchantQueryServiceClient
	merchantpb.MerchantCommandServiceClient
}

type MerchantDocumentClient struct {
	MerchantDocumentQueryClient   merchantdocumentpb.MerchantDocumentQueryServiceClient
	MerchantDocumentCommandClient merchantdocumentpb.MerchantDocumentCommandServiceClient
}

type OrderClient struct {
	orderpb.OrderQueryServiceClient
	orderpb.OrderCommandServiceClient
	orderpb.OrderStatsServiceClient
	orderpb.OrderStatsByMerchantServiceClient
}

type OrderItemClient struct {
	orderitempb.OrderItemQueryServiceClient
}

type ProductClient struct {
	productpb.ProductQueryServiceClient
	productpb.ProductCommandServiceClient
}

type TransactionClient struct {
	transactionpb.TransactionQueryServiceClient
	transactionpb.TransactionCommandServiceClient
	transactionpb.TransactionStatsServiceClient
	transactionpb.TransactionStatsByMerchantServiceClient
}

type AuthHandleGraphql struct {
	AuthClient authpb.AuthServiceClient
	Logger     logger.LoggerInterface
	Mapping    authgraphqlmapper.AuthGraphqlMapper
	Cache      auth_cache.AuthMencache
}

type RoleHandleGraphql struct {
	RoleClient RoleClient
	Logger     logger.LoggerInterface
	Mapping    rolegraphqlmapper.RoleGraphqlMapper
	Kafka      *kafka.Kafka
	Permission rolepermission.RolePermission
	Cache      role_cache.RoleMencache
}

type UserHandleGraphql struct {
	UserClient UserClient
	Logger     logger.LoggerInterface
	Mapping    usergraphqlmapper.UserGraphqlMapper
	Cache      user_cache.UserMencache
}

type CashierHandleGraphql struct {
	CashierClient CashierClient
	Logger        logger.LoggerInterface
	Mapping       cashiergraphqlmapper.CashierGraphqlMapper
	Cache         cashier_cache.CashierMencache
}

type CategoryHandleGraphql struct {
	CategoryClient CategoryClient
	Logger         logger.LoggerInterface
	Mapping        categorygraphqlmapper.CategoryGraphqlMapper
	Cache          category_cache.CategoryMencache
}

type MerchantHandleGraphql struct {
	MerchantClient MerchantClient
	Logger         logger.LoggerInterface
	Mapping        merchantgraphqlmapper.MerchantGraphqlMapper
	Cache          merchant_cache.MerchantMenCache
}

type MerchantDocumentHandleGraphql struct {
	MerchantClient MerchantDocumentClient
	Logger         logger.LoggerInterface
	Mapping        merchantdocumentgraphqlmapper.MerchantDocumentGraphqlMapper
	Cache          merchant_document_cache.MerchantDocumentMencache
}

type OrderHandleGraphql struct {
	OrderClient OrderClient
	Logger      logger.LoggerInterface
	Mapping     ordergraphqlmapper.OrderGraphqlMapper
	Cache       order_cache.OrderMencache
}

type OrderItemHandleGraphql struct {
	OrderItemClient OrderItemClient
	Logger          logger.LoggerInterface
	Mapping         orderitemgraphqlmapper.OrderItemGraphqlMapper
	Cache           orderitem_cache.OrderItemCache
}

type ProductHandleGraphql struct {
	ProductClient ProductClient
	Logger        logger.LoggerInterface
	Mapping       productgraphqlmapper.ProductGraphqlMapper
	Cache         product_cache.ProductMencache
	ImageUpload   upload_image.ImageUploads
}

type TransactionHandleGraphql struct {
	TransactionClient TransactionClient
	Logger            logger.LoggerInterface
	Mapping           transactiongraphqlmapper.TransactionGraphqlMapper
	Permission        merchantpermission.MerchantPermission
	Cache             transaction_cache.TransactionMencache
}

type ServiceConnections struct {
	AuthClient        *grpc.ClientConn
	CashierClient     *grpc.ClientConn
	CategoryClient    *grpc.ClientConn
	MerchantClient    *grpc.ClientConn
	OrderClient       *grpc.ClientConn
	OrderItemClient   *grpc.ClientConn
	ProductClient     *grpc.ClientConn
	RoleClient        *grpc.ClientConn
	TransactionClient *grpc.ClientConn
	UserClient        *grpc.ClientConn
}

type Deps struct {
	Clients  *ServiceConnections
	Logger   logger.LoggerInterface
	Kafka    *kafka.Kafka
	Mencache mencache.CacheApiGateway
}

func NewResolver(
	deps *Deps,
) *Resolver {
	observability, _ := observability.NewObservability(
		"graphql-client",
		deps.Logger,
	)

	resolverHandle := NewResolverHandler(observability, deps.Logger)

	store := deps.Mencache.GetStore()
	cacheAuth := auth_cache.NewMencache(store)
	cacheUser := user_cache.NewUserMencache(store)
	cacheRole := role_cache.NewRoleMencache(store)
	cacheMerchant := merchant_cache.NewMerchantMencache(store)
	cacheMerchantDocument := merchant_document_cache.NewMerchantDocumentMencache(store)
	cacheCashier := cashier_cache.NewCashierMencache(store)
	cacheCategory := category_cache.NewCategoryMencache(store)
	cacheOrder := order_cache.NewOrderMencache(store)
	cacheOrderItem := orderitem_cache.NewOrderItemCache(store)
	cacheProduct := product_cache.NewProductMencache(store)
	cacheTransaction := transaction_cache.NewTransactionMencache(store)

	return &Resolver{
		ResolverHandle: resolverHandle,
		AuthGraphql: AuthHandleGraphql{
			AuthClient: authpb.NewAuthServiceClient(deps.Clients.AuthClient),
			Logger:     deps.Logger,
			Mapping:    authgraphqlmapper.NewAuthGraphqlMapper(),
			Cache:      cacheAuth,
		},
		RoleGraphql: RoleHandleGraphql{
			RoleClient: RoleClient{
				RoleQueryServiceClient:   rolepb.NewRoleQueryServiceClient(deps.Clients.RoleClient),
				RoleCommandServiceClient: rolepb.NewRoleCommandServiceClient(deps.Clients.RoleClient),
			},
			Kafka:      deps.Kafka,
			Logger:     deps.Logger,
			Mapping:    rolegraphqlmapper.NewRoleGraphqlMapper(),
			Permission: rolepermission.NewRolePermission(deps.Kafka, "request-role", "response-role", 5*time.Second, deps.Logger, deps.Mencache),
			Cache:      cacheRole,
		},
		UserGraphql: UserHandleGraphql{
			UserClient: UserClient{
				UserQueryClient:   userpb.NewUserQueryServiceClient(deps.Clients.UserClient),
				UserCommandClient: userpb.NewUserCommandServiceClient(deps.Clients.UserClient),
			},
			Logger:  deps.Logger,
			Mapping: usergraphqlmapper.NewUserGraphqlMapper(),
			Cache:   cacheUser,
		},
		CashierGraphql: CashierHandleGraphql{
			CashierClient: CashierClient{
				CashierQueryServiceClient:           cashierpb.NewCashierQueryServiceClient(deps.Clients.CashierClient),
				CashierCommandServiceClient:         cashierpb.NewCashierCommandServiceClient(deps.Clients.CashierClient),
				CashierStatsServiceClient:           cashierpb.NewCashierStatsServiceClient(deps.Clients.CashierClient),
				CashierStatsByIdServiceClient:       cashierpb.NewCashierStatsByIdServiceClient(deps.Clients.CashierClient),
				CashierStatsByMerchantServiceClient: cashierpb.NewCashierStatsByMerchantServiceClient(deps.Clients.CashierClient),
			},
			Logger:  deps.Logger,
			Mapping: cashiergraphqlmapper.NewCashierGraphqlMapper(),
			Cache:   cacheCashier,
		},
		CategoryGraphql: CategoryHandleGraphql{
			CategoryClient: CategoryClient{
				CategoryQueryServiceClient:           categorypb.NewCategoryQueryServiceClient(deps.Clients.CategoryClient),
				CategoryCommandServiceClient:         categorypb.NewCategoryCommandServiceClient(deps.Clients.CategoryClient),
				CategoryStatsServiceClient:           categorypb.NewCategoryStatsServiceClient(deps.Clients.CategoryClient),
				CategoryStatsByIdServiceClient:       categorypb.NewCategoryStatsByIdServiceClient(deps.Clients.CategoryClient),
				CategoryStatsByMerchantServiceClient: categorypb.NewCategoryStatsByMerchantServiceClient(deps.Clients.CategoryClient),
			},
			Logger:  deps.Logger,
			Mapping: categorygraphqlmapper.NewCategoryGraphqlMapper(),
			Cache:   cacheCategory,
		},
		MerchantGraphql: MerchantHandleGraphql{
			MerchantClient: MerchantClient{
				MerchantQueryServiceClient:   merchantpb.NewMerchantQueryServiceClient(deps.Clients.MerchantClient),
				MerchantCommandServiceClient: merchantpb.NewMerchantCommandServiceClient(deps.Clients.MerchantClient),
			},
			Logger:  deps.Logger,
			Mapping: merchantgraphqlmapper.NewMerchantGraphqlMapper(),
			Cache:   cacheMerchant,
		},
		MerchantDocumentGraphql: MerchantDocumentHandleGraphql{
			MerchantClient: MerchantDocumentClient{
				MerchantDocumentQueryClient:   merchantdocumentpb.NewMerchantDocumentQueryServiceClient(deps.Clients.MerchantClient),
				MerchantDocumentCommandClient: merchantdocumentpb.NewMerchantDocumentCommandServiceClient(deps.Clients.MerchantClient),
			},
			Logger:  deps.Logger,
			Mapping: merchantdocumentgraphqlmapper.NewMerchantDocumentGraphqlMapper(),
			Cache:   cacheMerchantDocument,
		},
		OrderGraphql: OrderHandleGraphql{
			OrderClient: OrderClient{
				OrderQueryServiceClient:           orderpb.NewOrderQueryServiceClient(deps.Clients.OrderClient),
				OrderCommandServiceClient:         orderpb.NewOrderCommandServiceClient(deps.Clients.OrderClient),
				OrderStatsServiceClient:           orderpb.NewOrderStatsServiceClient(deps.Clients.OrderClient),
				OrderStatsByMerchantServiceClient: orderpb.NewOrderStatsByMerchantServiceClient(deps.Clients.OrderClient),
			},
			Logger:  deps.Logger,
			Mapping: ordergraphqlmapper.NewOrderGraphqlMapper(),
			Cache:   cacheOrder,
		},
		OrderItemGraphql: OrderItemHandleGraphql{
			OrderItemClient: OrderItemClient{
				OrderItemQueryServiceClient: orderitempb.NewOrderItemQueryServiceClient(deps.Clients.OrderItemClient),
			},
			Logger:  deps.Logger,
			Mapping: orderitemgraphqlmapper.NewOrderItemGraphqlMapper(),
			Cache:   cacheOrderItem,
		},
		ProductGraphql: ProductHandleGraphql{
			ProductClient: ProductClient{
				ProductQueryServiceClient:   productpb.NewProductQueryServiceClient(deps.Clients.ProductClient),
				ProductCommandServiceClient: productpb.NewProductCommandServiceClient(deps.Clients.ProductClient),
			},
			Logger:      deps.Logger,
			Mapping:     productgraphqlmapper.NewProductGraphqlMapper(),
			Cache:       cacheProduct,
			ImageUpload: upload_image.NewImageUpload(deps.Logger),
		},
		TransactionGraphql: TransactionHandleGraphql{
			TransactionClient: TransactionClient{
				TransactionQueryServiceClient:           transactionpb.NewTransactionQueryServiceClient(deps.Clients.TransactionClient),
				TransactionCommandServiceClient:         transactionpb.NewTransactionCommandServiceClient(deps.Clients.TransactionClient),
				TransactionStatsServiceClient:           transactionpb.NewTransactionStatsServiceClient(deps.Clients.TransactionClient),
				TransactionStatsByMerchantServiceClient: transactionpb.NewTransactionStatsByMerchantServiceClient(deps.Clients.TransactionClient),
			},
			Logger:     deps.Logger,
			Mapping:    transactiongraphqlmapper.NewTransactionGraphqlMapper(),
			Permission: merchantpermission.NewMerchantPermission(deps.Kafka, "request-transaction", "response-transaction", 5*time.Second, deps.Logger),
			Cache:      cacheTransaction,
		},
	}
}

func (h *Resolver) handleGraphQLError(err error, operation string) *errors.AppError {
	if err == nil {
		return nil
	}

	var appErr *errors.AppError
	if errorstd.As(err, &appErr) {
		return appErr
	}

	return errors.NewInternalError(err).WithMessage("Failed to " + operation)
}

func (r *Resolver) parseValidationErrors(err error) []sharedErrors.ValidationError {
	var validationErrs []sharedErrors.ValidationError

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			validationErrs = append(validationErrs, sharedErrors.ValidationError{
				Field:   fe.Field(),
				Message: r.getValidationMessage(fe),
			})
		}
		return validationErrs
	}

	return []sharedErrors.ValidationError{
		{
			Field:   "general",
			Message: err.Error(),
		},
	}
}

func (r *Resolver) getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Must be at least %s", fe.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s", fe.Param())
	case "gte":
		return fmt.Sprintf("Must be greater than or equal to %s", fe.Param())
	case "lte":
		return fmt.Sprintf("Must be less than or equal to %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", fe.Param())
	default:
		return fmt.Sprintf("Validation failed on '%s' tag", fe.Tag())
	}
}
