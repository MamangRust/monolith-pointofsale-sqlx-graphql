package graphtest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	mycontext "github.com/MamangRust/monolith-pointofsale-apigateway/internal/context"
	graph "github.com/MamangRust/monolith-pointofsale-apigateway/internal/handler"
	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/middlewares"
	mencache "github.com/MamangRust/monolith-pointofsale-apigateway/internal/redis"
	"github.com/MamangRust/monolith-pointofsale-pkg/auth"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/redis/go-redis/v9"
)

// WithUserID wraps the internal mycontext.WithUserID function to make it publicly accessible.
func WithUserID(ctx context.Context, userID int) context.Context {
	return mycontext.WithUserID(ctx, userID)
}

// Resolver is an alias for the internal graph.Resolver type.
type Resolver = graph.Resolver

// ServiceConnections is an alias for the internal ServiceConnections type.
type ServiceConnections = graph.ServiceConnections

// NewResolver creates a new Resolver from deps.
func NewResolver(conns *ServiceConnections, log logger.LoggerInterface, redisClient *redis.Client) *Resolver {
	store := mencache.NewCacheApiGateway(&mencache.Deps{
		Redis:  redisClient,
		Logger: log,
	})
	return graph.NewResolver(&graph.Deps{
		Clients:  conns,
		Logger:   log,
		Mencache: store,
		Kafka:    nil,
	})
}

// NewHandler creates a GraphQL HTTP handler from a resolver.
func NewHandler(resolver *Resolver) *handler.Server {
	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
	}))
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	return srv
}

// NewAuthHandler creates a GraphQL HTTP handler wrapped with the apigateway
// auth middleware (public ops like registerUser/loginUser bypass auth; the
// rest require a valid Bearer token).
func NewAuthHandler(resolver *Resolver, tokenManager auth.TokenManager, log logger.LoggerInterface) http.Handler {
	return middlewares.AuthMiddleware(tokenManager, log)(NewHandler(resolver))
}

// GraphQLQuery represents a GraphQL request payload.
type GraphQLQuery struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// GraphQLResponse wraps the standard GraphQL response structure.
type GraphQLResponse struct {
	Data   map[string]interface{} `json:"data"`
	Errors []GraphQLError         `json:"errors,omitempty"`
}

// GraphQLError represents a single error in the GraphQL response.
type GraphQLError struct {
	Message string `json:"message"`
}

// ExecuteGraphQL executes a query against the handler.
func ExecuteGraphQL(srv http.Handler, query string, variables map[string]interface{}, authToken string) (*GraphQLResponse, error) {
	gqlReq := GraphQLQuery{
		Query:     query,
		Variables: variables,
	}
	body, err := json.Marshal(gqlReq)
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest(http.MethodPost, "/query", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	respBody, err := io.ReadAll(rec.Result().Body)
	if err != nil {
		return nil, err
	}

	var resp GraphQLResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
