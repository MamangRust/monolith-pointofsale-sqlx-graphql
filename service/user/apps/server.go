package apps

import (
	"time"

	pbrole "github.com/MamangRust/monolith-pointofsale-pb/role"
	pbuser "github.com/MamangRust/monolith-pointofsale-pb/user"
	pbuserrole "github.com/MamangRust/monolith-pointofsale-pb/user_role"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter"
	"github.com/MamangRust/monolith-pointofsale-pkg/hash"
	"github.com/MamangRust/monolith-pointofsale-pkg/resilience"
	"github.com/MamangRust/monolith-pointofsale-pkg/server"

	"github.com/MamangRust/monolith-pointofsale-shared/cache"
	"github.com/MamangRust/monolith-pointofsale-shared/convert"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"

	mencache "github.com/MamangRust/monolith-pointofsale-user/cache"
	"github.com/MamangRust/monolith-pointofsale-user/handler"
	"github.com/MamangRust/monolith-pointofsale-user/repository"
	"github.com/MamangRust/monolith-pointofsale-user/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	roleAddr := convert.EnvOr("GRPC_ROLE_ADDR", "localhost:50052")

	srv.Logger.Info("Connecting to gRPC microservices from User microservice",
		zap.String("role", roleAddr),
	)

	roleConn, err := server.NewGRPCClient(roleAddr, resilience.GuardDialOption("role", srv.Logger))
	if err != nil {
		return nil, err
	}

	go func() {
		<-srv.Ctx.Done()
		srv.Logger.Info("Closing gRPC client connections in User microservice")
		roleConn.Close()
	}()

	roleQueryClient := pbrole.NewRoleQueryServiceClient(roleConn)
	userRoleClient := pbuserrole.NewUserRoleServiceClient(roleConn)

	guardRole := resilience.NewDependencyGuard("role", 5, 30, 100, 3*time.Second, srv.Logger)
	guardUserRole := resilience.NewDependencyGuard("user_role", 5, 30, 100, 3*time.Second, srv.Logger)

	repos := repository.NewRepositories(&repository.Deps{
		Db:              srv.DB,
		RoleQueryClient: roleQueryClient,
		UserRoleClient:  userRoleClient,
		Guards: repository.GuardOptions{
			Role:     []adapter.GuardOption{adapter.WithDependencyGuard(guardRole)},
			UserRole: []adapter.GuardOption{adapter.WithDependencyGuard(guardUserRole)},
		},
	})
	hash := hash.NewHashingPassword()
	traceLoggerObservability := observability.NewTraceLoggerObservability(srv.Logger)

	cacheMetrics, err := observability.NewCacheMetrics("user")
	if err != nil {
		return nil, err
	}

	cacheStore := cache.NewCacheStore(srv.Redis, srv.Logger, cacheMetrics)
	mencacheObj := mencache.NewMencache(cacheStore)

	services := service.NewService(&service.Deps{
		Mencache:      mencacheObj,
		Repositories:  repos,
		Hash:          hash,
		Logger:        srv.Logger,
		Observability: traceLoggerObservability,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbuser.RegisterUserQueryServiceServer(gs, handlers.User)
		pbuser.RegisterUserCommandServiceServer(gs, handlers.UserCommand)
	}

	return srv, nil
}
