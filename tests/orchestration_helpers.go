package tests

import (
	"bytes"
	"context"
	"mime/multipart"

	"github.com/MamangRust/monolith-pointofsale-pkg/auth"
	"github.com/MamangRust/monolith-pointofsale-pkg/hash"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Role
	role_cache "github.com/MamangRust/monolith-pointofsale-role/cache"
	role_handler "github.com/MamangRust/monolith-pointofsale-role/handler"
	role_repo "github.com/MamangRust/monolith-pointofsale-role/repository"
	role_service "github.com/MamangRust/monolith-pointofsale-role/service"

	// User
	user_cache "github.com/MamangRust/monolith-pointofsale-user/cache"
	user_handler "github.com/MamangRust/monolith-pointofsale-user/handler"
	user_repo "github.com/MamangRust/monolith-pointofsale-user/repository"
	user_service "github.com/MamangRust/monolith-pointofsale-user/service"

	// Auth
	auth_cache "github.com/MamangRust/monolith-pointofsale-auth/cache"
	auth_handler "github.com/MamangRust/monolith-pointofsale-auth/handler"
	auth_repo "github.com/MamangRust/monolith-pointofsale-auth/repository"
	auth_service "github.com/MamangRust/monolith-pointofsale-auth/service"

	// Category
	category_cache "github.com/MamangRust/monolith-pointofsale-category/cache"
	category_handler "github.com/MamangRust/monolith-pointofsale-category/handler"
	category_repo "github.com/MamangRust/monolith-pointofsale-category/repository"
	category_service "github.com/MamangRust/monolith-pointofsale-category/service"

	// Product
	product_cache "github.com/MamangRust/monolith-pointofsale-product/cache"
	product_handler "github.com/MamangRust/monolith-pointofsale-product/handler"
	product_repo "github.com/MamangRust/monolith-pointofsale-product/repository"
	product_service "github.com/MamangRust/monolith-pointofsale-product/service"

	// Merchant
	merchant_cache "github.com/MamangRust/monolith-pointofsale-merchant/cache"
	merchant_handler "github.com/MamangRust/monolith-pointofsale-merchant/handler"
	merchant_repo "github.com/MamangRust/monolith-pointofsale-merchant/repository"
	merchant_service "github.com/MamangRust/monolith-pointofsale-merchant/service"

	// Order
	order_cache "github.com/MamangRust/monolith-pointofsale-order/cache"
	order_handler "github.com/MamangRust/monolith-pointofsale-order/handler"
	order_repo "github.com/MamangRust/monolith-pointofsale-order/repository"
	order_service "github.com/MamangRust/monolith-pointofsale-order/service"

	// Transaction
	transaction_cache "github.com/MamangRust/monolith-pointofsale-transacton/cache"
	transaction_handler "github.com/MamangRust/monolith-pointofsale-transacton/handler"
	transaction_repo "github.com/MamangRust/monolith-pointofsale-transacton/repository"
	transaction_service "github.com/MamangRust/monolith-pointofsale-transacton/service"

	// Order Item
	order_item_cache "github.com/MamangRust/monolith-pointofsale-order-item/cache"
	order_item_handler "github.com/MamangRust/monolith-pointofsale-order-item/handler"
	order_item_repo "github.com/MamangRust/monolith-pointofsale-order-item/repository"
	order_item_service "github.com/MamangRust/monolith-pointofsale-order-item/service"

	// Cashier
	mencache "github.com/MamangRust/monolith-pointofsale-cashier/cache"
	cashier_handler "github.com/MamangRust/monolith-pointofsale-cashier/handler"
	cashier_repo "github.com/MamangRust/monolith-pointofsale-cashier/repository"
	cashier_service "github.com/MamangRust/monolith-pointofsale-cashier/service"
	pb "github.com/MamangRust/monolith-pointofsale-pb"
	pbcashier "github.com/MamangRust/monolith-pointofsale-pb/cashier"
	pbcategory "github.com/MamangRust/monolith-pointofsale-pb/category"
	pbmerchant "github.com/MamangRust/monolith-pointofsale-pb/merchant"
	pbmerchant_document "github.com/MamangRust/monolith-pointofsale-pb/merchant_document"
	pborder "github.com/MamangRust/monolith-pointofsale-pb/order"
	pbproduct "github.com/MamangRust/monolith-pointofsale-pb/product"
	pbrole "github.com/MamangRust/monolith-pointofsale-pb/role"
	pbtransaction "github.com/MamangRust/monolith-pointofsale-pb/transaction"
	pbuser "github.com/MamangRust/monolith-pointofsale-pb/user"
	pbuserrole "github.com/MamangRust/monolith-pointofsale-pb/user_role"
)

func (s *BaseTestSuite) SetupRoleService() {
	cacheStore := s.GetCacheStore()
	queries := s.SQLxDB()

	roleMencache := role_cache.NewMencache(cacheStore)
	roleRepos := role_repo.NewRepositories(queries)
	roleSvc := role_service.NewService(&role_service.Deps{
		Repositories:  roleRepos,
		Logger:        s.Log,
		Mencache:      roleMencache,
		Observability: s.Obs,
	})
	roleGapi := role_handler.NewHandler(&role_handler.Deps{
		Service: roleSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pbrole.RegisterRoleQueryServiceServer(server, roleGapi.Role)
	pbrole.RegisterRoleCommandServiceServer(server, roleGapi.RoleCommand)
	pbuserrole.RegisterUserRoleServiceServer(server, roleGapi.UserRole)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.Conns["role"] = conn
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupUserService() {
	if _, ok := s.Conns["role"]; !ok {
		s.SetupRoleService()
	}

	cacheStore := s.GetCacheStore()
	queries := s.SQLxDB()
	hasher := hash.NewHashingPassword()

	userMencache := user_cache.NewMencache(cacheStore)
	userRepos := user_repo.NewRepositories(&user_repo.Deps{
		Db:              queries,
		RoleQueryClient: pbrole.NewRoleQueryServiceClient(s.Conns["role"]),
		UserRoleClient:  pbuserrole.NewUserRoleServiceClient(s.Conns["role"]),
	})
	userSvc := user_service.NewService(&user_service.Deps{
		Repositories:  userRepos,
		Logger:        s.Log,
		Hash:          hasher,
		Mencache:      userMencache,
		Observability: s.Obs,
	})
	userGapi := user_handler.NewHandler(&user_handler.Deps{
		Service: userSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pbuser.RegisterUserQueryServiceServer(server, userGapi.User)
	pbuser.RegisterUserCommandServiceServer(server, userGapi.UserCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.Conns["user"] = conn
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupAuthService() {
	if _, ok := s.Conns["user"]; !ok {
		s.SetupUserService()
	}

	cacheStore := s.GetCacheStore()
	queries := s.SQLxDB()
	hasher := hash.NewHashingPassword()
	tokenManager, _ := auth.NewManager("mysecret")

	authRepos := auth_repo.NewRepositories(&auth_repo.Deps{
		Db:                queries,
		UserQueryClient:   pbuser.NewUserQueryServiceClient(s.Conns["user"]),
		UserCommandClient: pbuser.NewUserCommandServiceClient(s.Conns["user"]),
		RoleQueryClient:   pbrole.NewRoleQueryServiceClient(s.Conns["role"]),
		RoleCommandClient: pbrole.NewRoleCommandServiceClient(s.Conns["role"]),
		UserRoleClient:    pbuserrole.NewUserRoleServiceClient(s.Conns["role"]),
	})
	authMencache := auth_cache.NewMencache(cacheStore)
	authSvc := auth_service.NewService(&auth_service.Deps{
		Repositories:  authRepos,
		Logger:        s.Log,
		Mencache:      authMencache,
		Token:         tokenManager,
		Hash:          hasher,
		Kafka:         nil,
		Observability: s.Obs,
	})
	authGapi := auth_handler.NewAuthHandleGrpc(authSvc, s.Log)
	server := grpc.NewServer()
	pb.RegisterAuthServiceServer(server, authGapi)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.Conns["auth"] = conn
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupCategoryService() {
	cacheStore := s.GetCacheStore()
	queries := s.SQLxDB()

	catMencache := category_cache.NewMencache(cacheStore)
	catRepos := category_repo.NewRepositories(queries)
	catSvc := category_service.NewService(&category_service.Deps{
		Mencache:      catMencache,
		Repositories:  catRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	catGapi := category_handler.NewHandler(&category_handler.Deps{
		Service: catSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pbcategory.RegisterCategoryQueryServiceServer(server, catGapi.Category)
	pbcategory.RegisterCategoryCommandServiceServer(server, catGapi.CategoryCommand)
	pbcategory.RegisterCategoryStatsServiceServer(server, catGapi.CategoryStats)
	pbcategory.RegisterCategoryStatsByIdServiceServer(server, catGapi.CategoryStatsById)
	pbcategory.RegisterCategoryStatsByMerchantServiceServer(server, catGapi.CategoryStatsByMerchant)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["category"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupProductService() {
	if _, ok := s.Conns["category"]; !ok {
		s.SetupCategoryService()
	}
	if _, ok := s.Conns["merchant"]; !ok {
		s.SetupMerchantService()
	}

	cacheStore := s.GetCacheStore()
	queries := s.SQLxDB()

	prodMencache := product_cache.NewMencache(cacheStore)
	prodRepos := product_repo.NewRepositories(queries,
		pbcategory.NewCategoryQueryServiceClient(s.Conns["category"]),
		pbmerchant.NewMerchantQueryServiceClient(s.Conns["merchant"]),
	)
	prodSvc := product_service.NewService(&product_service.Deps{
		Mencache:      prodMencache,
		Repositories:  prodRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	prodGapi := product_handler.NewHandler(&product_handler.Deps{
		Service: prodSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pbproduct.RegisterProductQueryServiceServer(server, prodGapi.Product)
	pbproduct.RegisterProductCommandServiceServer(server, prodGapi.ProductCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["product"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupMerchantService() {
	if _, ok := s.Conns["user"]; !ok {
		s.SetupUserService()
	}

	cacheStore := s.GetCacheStore()
	queries := s.SQLxDB()

	merchantMencache := merchant_cache.NewMencache(cacheStore)
	userQueryClient := NewUserClient(s.Conns["user"])
	merchantRepos := merchant_repo.NewRepositories(queries, userQueryClient)
	merchantSvc := merchant_service.NewService(&merchant_service.Deps{
		Mencache:      merchantMencache,
		Repositories:  merchantRepos,
		Logger:        s.Log,
		Observability: s.Obs,
		Kafka:         nil,
	})
	merchantGapi := merchant_handler.NewHandler(&merchant_handler.Deps{
		Service: merchantSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pbmerchant.RegisterMerchantQueryServiceServer(server, merchantGapi.Merchant)
	pbmerchant.RegisterMerchantCommandServiceServer(server, merchantGapi.MerchantCommand)
	pbmerchant_document.RegisterMerchantDocumentQueryServiceServer(server, merchantGapi.MerchantDocument)
	pbmerchant_document.RegisterMerchantDocumentCommandServiceServer(server, merchantGapi.MerchantDocumentCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["merchant"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupOrderService() {
	if _, ok := s.Conns["cashier"]; !ok {
		s.SetupCashierService()
	}
	if _, ok := s.Conns["merchant"]; !ok {
		s.SetupMerchantService()
	}
	if _, ok := s.Conns["product"]; !ok {
		s.SetupProductService()
	}
	if _, ok := s.Conns["order-item"]; !ok {
		s.SetupOrderItemService()
	}

	cacheStore := s.GetCacheStore()
	queries := s.SQLxDB()

	cashierClient := NewCashierClient(s.Conns["cashier"])
	merchantClient := NewMerchantClient(s.Conns["merchant"])
	productClient := NewProductClient(s.Conns["product"])
	orderItemClient := NewOrderItemClient(s.Conns["order-item"])

	orderMencache := order_cache.NewMencache(cacheStore)
	orderRepos := order_repo.NewRepositories(&order_repo.Deps{
		Db:                     queries,
		CashierQueryClient:     cashierClient,
		MerchantQueryClient:    merchantClient,
		ProductQueryClient:     productClient,
		ProductCommandClient:   productClient,
		OrderItemQueryClient:   orderItemClient,
		OrderItemCommandClient: orderItemClient,
	})
	orderSvc := order_service.NewService(&order_service.Deps{
		Mencache:      orderMencache,
		Repositories:  orderRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	orderGapi := order_handler.NewHandler(&order_handler.Deps{
		Service: orderSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pborder.RegisterOrderQueryServiceServer(server, orderGapi.Order)
	pborder.RegisterOrderCommandServiceServer(server, orderGapi.OrderCommand)
	pborder.RegisterOrderStatsServiceServer(server, orderGapi.OrderStats)
	pborder.RegisterOrderStatsByMerchantServiceServer(server, orderGapi.OrderStatsByMerchant)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["order"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupTransactionService() {
	if _, ok := s.Conns["cashier"]; !ok {
		s.SetupCashierService()
	}
	if _, ok := s.Conns["merchant"]; !ok {
		s.SetupMerchantService()
	}
	if _, ok := s.Conns["order"]; !ok {
		s.SetupOrderService()
	}
	if _, ok := s.Conns["order-item"]; !ok {
		s.SetupOrderItemService()
	}

	cacheStore := s.GetCacheStore()
	queries := s.SQLxDB()

	cashierClient := NewCashierClient(s.Conns["cashier"])
	merchantClient := NewMerchantClient(s.Conns["merchant"])
	orderClient := NewOrderClient(s.Conns["order"])
	orderItemClient := NewOrderItemClient(s.Conns["order-item"])

	transactionMencache := transaction_cache.NewMencache(cacheStore)
	transactionRepos := transaction_repo.NewRepositories(queries, cashierClient, merchantClient, orderClient, orderItemClient)
	transactionSvc := transaction_service.NewService(&transaction_service.Deps{
		Mencache:      transactionMencache,
		Repositories:  transactionRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	transactionGapi := transaction_handler.NewHandler(&transaction_handler.Deps{
		Service: transactionSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pbtransaction.RegisterTransactionQueryServiceServer(server, transactionGapi.Transaction)
	pbtransaction.RegisterTransactionCommandServiceServer(server, transactionGapi.TransactionCommand)
	pbtransaction.RegisterTransactionStatsServiceServer(server, transactionGapi.TransactionStats)
	pbtransaction.RegisterTransactionStatsByMerchantServiceServer(server, transactionGapi.TransactionStatsByMerchant)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["transaction"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupOrderItemService() {
	cacheStore := s.GetCacheStore()
	queries := s.SQLxDB()

	itemMencache := order_item_cache.NewMencache(cacheStore)
	itemRepos := order_item_repo.NewRepositories(queries)
	itemSvc := order_item_service.NewService(&order_item_service.Deps{
		Mencache:      itemMencache,
		Repositories:  itemRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	itemGapi := order_item_handler.NewHandler(&order_item_handler.Deps{
		Service: itemSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterOrderItemQueryServiceServer(server, itemGapi.OrderItem)
	pb.RegisterOrderItemCommandServiceServer(server, itemGapi.OrderItemCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["order-item"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupCashierService() {
	if _, ok := s.Conns["user"]; !ok {
		s.SetupUserService()
	}
	if _, ok := s.Conns["merchant"]; !ok {
		s.SetupMerchantService()
	}

	cacheStore := s.GetCacheStore()
	queries := s.SQLxDB()

	cashierMencache := mencache.NewMencache(cacheStore)
	userClient := NewUserClient(s.Conns["user"])
	merchantClient := NewMerchantClient(s.Conns["merchant"])
	cashierRepos := cashier_repo.NewRepositories(queries, userClient, merchantClient)
	cashierSvc := cashier_service.NewService(&cashier_service.Deps{
		Ctx:           context.Background(),
		Mencache:      cashierMencache,
		Repositories:  cashierRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	cashierGapi := cashier_handler.NewHandler(&cashier_handler.Deps{
		Service: cashierSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pbcashier.RegisterCashierQueryServiceServer(server, cashierGapi.Cashier)
	pbcashier.RegisterCashierCommandServiceServer(server, cashierGapi.CashierCommand)
	pbcashier.RegisterCashierStatsServiceServer(server, cashierGapi.CashierStats)
	pbcashier.RegisterCashierStatsByIdServiceServer(server, cashierGapi.CashierStatsById)
	pbcashier.RegisterCashierStatsByMerchantServiceServer(server, cashierGapi.CashierStatsByMerchant)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["cashier"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) dial(addr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	return conn
}

func (s *BaseTestSuite) GetCacheStore() *cache.CacheStore {
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	return cache.NewCacheStore(s.ts.RedisClient(), s.Log, cacheMetrics)
}

func (s *BaseTestSuite) BuildMultipartRequestBody(fields map[string]string, fieldName, fileName string) ([]byte, string) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range fields {
		fw, _ := w.CreateFormField(key)
		fw.Write([]byte(r))
	}
	fw, _ := w.CreateFormFile(fieldName, fileName)
	fw.Write([]byte("dummy image content"))
	w.Close()
	return b.Bytes(), w.FormDataContentType()
}

// StartRoleUserServices boots role-service (including UserRoleService) and
// user-service as in-process gRPC servers. Auth reaches both over gRPC, so auth
// integration tests must start them before wiring the auth repositories. The
// caller owns cleanup of the returned connections and servers.
func StartRoleUserServices(
	queries *sqlx.DB,
	log logger.LoggerInterface,
	cacheStore *cache.CacheStore,
	obs observability.TraceLoggerObservability,
) (roleConn, userConn *grpc.ClientConn, servers []*grpc.Server, err error) {
	roleMencache := role_cache.NewMencache(cacheStore)
	roleRepos := role_repo.NewRepositories(queries)
	roleSvc := role_service.NewService(&role_service.Deps{
		Repositories:  roleRepos,
		Logger:        log,
		Mencache:      roleMencache,
		Observability: obs,
	})
	roleGapi := role_handler.NewHandler(&role_handler.Deps{
		Service: roleSvc,
		Logger:  log,
	})
	roleServer := grpc.NewServer()
	pbrole.RegisterRoleQueryServiceServer(roleServer, roleGapi.Role)
	pbrole.RegisterRoleCommandServiceServer(roleServer, roleGapi.RoleCommand)
	pbuserrole.RegisterUserRoleServiceServer(roleServer, roleGapi.UserRole)

	roleAddr, err := RunGRPCServer(roleServer)
	if err != nil {
		return nil, nil, nil, err
	}
	roleConn, err = grpc.NewClient(roleAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		roleServer.Stop()
		return nil, nil, nil, err
	}

	userMencache := user_cache.NewMencache(cacheStore)
	userRepos := user_repo.NewRepositories(&user_repo.Deps{
		Db:              queries,
		RoleQueryClient: pbrole.NewRoleQueryServiceClient(roleConn),
		UserRoleClient:  pbuserrole.NewUserRoleServiceClient(roleConn),
	})
	userSvc := user_service.NewService(&user_service.Deps{
		Repositories:  userRepos,
		Logger:        log,
		Hash:          hash.NewHashingPassword(),
		Mencache:      userMencache,
		Observability: obs,
	})
	userGapi := user_handler.NewHandler(&user_handler.Deps{
		Service: userSvc,
		Logger:  log,
	})
	userServer := grpc.NewServer()
	pbuser.RegisterUserQueryServiceServer(userServer, userGapi.User)
	pbuser.RegisterUserCommandServiceServer(userServer, userGapi.UserCommand)

	userAddr, err := RunGRPCServer(userServer)
	if err != nil {
		roleConn.Close()
		roleServer.Stop()
		return nil, nil, nil, err
	}
	userConn, err = grpc.NewClient(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		userServer.Stop()
		roleConn.Close()
		roleServer.Stop()
		return nil, nil, nil, err
	}

	return roleConn, userConn, []*grpc.Server{roleServer, userServer}, nil
}
