package apps

import (
	"fmt"
	"os"
	"time"

	pb "github.com/MamangRust/monolith-pointofsale-pb"

	mencache "github.com/MamangRust/monolith-pointofsale-auth/cache"
	"github.com/MamangRust/monolith-pointofsale-auth/handler"
	"github.com/MamangRust/monolith-pointofsale-auth/repository"
	"github.com/MamangRust/monolith-pointofsale-auth/service"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter"
	"github.com/MamangRust/monolith-pointofsale-pkg/auth"
	"github.com/MamangRust/monolith-pointofsale-pkg/hash"
	"github.com/MamangRust/monolith-pointofsale-pkg/kafka"
	"github.com/MamangRust/monolith-pointofsale-pkg/outbox"
	"github.com/MamangRust/monolith-pointofsale-pkg/resilience"
	"github.com/MamangRust/monolith-pointofsale-pkg/server"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"

	pbrole "github.com/MamangRust/monolith-pointofsale-pb/role"
	pbuser "github.com/MamangRust/monolith-pointofsale-pb/user"
	pbuserrole "github.com/MamangRust/monolith-pointofsale-pb/user_role"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	tokenManager, err := auth.NewManager(viper.GetString("SECRET_KEY"))
	if err != nil {
		return nil, err
	}

	roleAddr := viper.GetString("GRPC_ROLE_ADDR")

	userAddr := viper.GetString("GRPC_USER_ADDR")

	roleConn, err := grpc.NewClient(
		roleAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to role service: %w", err)
	}

	userConn, err := grpc.NewClient(
		userAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	roleQueryClient := pbrole.NewRoleQueryServiceClient(roleConn)
	roleCommandClient := pbrole.NewRoleCommandServiceClient(roleConn)
	userQueryClient := pbuser.NewUserQueryServiceClient(userConn)
	userCommandClient := pbuser.NewUserCommandServiceClient(userConn)
	userRoleClient := pbuserrole.NewUserRoleServiceClient(roleConn)

	guardUser := resilience.NewDependencyGuard("user", 5, 30, 100, 3*time.Second, srv.Logger)
	guardRole := resilience.NewDependencyGuard("role", 5, 30, 100, 3*time.Second, srv.Logger)
	guardUserRole := resilience.NewDependencyGuard("user_role", 5, 30, 100, 3*time.Second, srv.Logger)

	hasher := hash.NewHashingPassword()

	repositories := repository.NewRepositories(&repository.Deps{
		Db:                srv.DB,
		UserQueryClient:   userQueryClient,
		UserCommandClient: userCommandClient,
		RoleQueryClient:   roleQueryClient,
		RoleCommandClient: roleCommandClient,
		UserRoleClient:    userRoleClient,
		Guard: repository.GuardOptions{
			User: []adapter.GuardOption{
				adapter.WithDependencyGuard(guardUser),
			},
			Role: []adapter.GuardOption{
				adapter.WithDependencyGuard(guardRole),
			},
			UserRole: []adapter.GuardOption{
				adapter.WithDependencyGuard(guardUserRole),
			},
		},
	})

	var myKafka *kafka.Kafka
	if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		myKafka = kafka.NewKafka(srv.Logger, []string{brokers})
	}

	mencacheObj := mencache.NewMencache(srv.CacheStore)

	outboxService := outbox.NewOutboxService(srv.DB, myKafka, srv.Logger)

	services := service.NewService(&service.Deps{
		Mencache:      mencacheObj,
		Repositories:  repositories,
		Hash:          hasher,
		Token:         tokenManager,
		Logger:        srv.Logger,
		Kafka:         myKafka,
		Pool:          srv.DB,
		Outbox:        outboxService,
		Observability: observability.NewTraceLoggerObservability(srv.Logger),
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterAuthServiceServer(gs, handlers.Auth)
	}

	// Start the outbox relay so events committed with the business writes are
	// published to Kafka with durable retry and dead-letter semantics.
	go outboxService.Start(srv.Ctx, outbox.OutboxRelayInterval, outbox.OutboxRelayBatchSize)

	return srv, nil
}
