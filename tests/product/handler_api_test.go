package product_test

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/MamangRust/monolith-pointofsale-pkg/auth"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	tests "github.com/MamangRust/monolith-pointofsale-test"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"

	graphtest "github.com/MamangRust/monolith-pointofsale-apigateway/graphtest"
)

type ProductApiTestSuite struct {
	tests.BaseTestSuite
	redisClient  *redis.Client
	handler      http.Handler
	tokenMgr     auth.TokenManager
	token        string
	productID    int
	merchantID   int
	categoryID   int
	categoryName string
}

func (s *ProductApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()

	ctx := context.Background()
	userID := s.SeedUser(ctx)
	s.categoryID = s.SeedCategory(ctx)
	s.merchantID = s.SeedMerchant(ctx, userID)
	s.categoryName = "Seed Category"

	s.redisClient = s.RedisClient()
	s.redisClient.FlushAll(ctx)

	log, _ := logger.NewLogger("test", nil)

	tokenMgr, err := auth.NewManager("mysecret")
	s.Require().NoError(err)
	s.tokenMgr = tokenMgr
	s.token, err = tokenMgr.GenerateToken(userID, "api")
	s.Require().NoError(err)

	conns := &graphtest.ServiceConnections{
		UserClient:     s.Conns["user"],
		CategoryClient: s.Conns["category"],
		MerchantClient: s.Conns["merchant"],
		ProductClient:  s.Conns["product"],
	}
	resolver := graphtest.NewResolver(conns, log, s.redisClient)
	s.handler = graphtest.NewAuthHandler(resolver, tokenMgr, log)
}

func (s *ProductApiTestSuite) TearDownSuite() {
	if s.redisClient != nil {
		s.redisClient.Close()
	}
	s.BaseTestSuite.TearDownSuite()
}

// executeGraphQLUpload sends a multipart GraphQL request with a file upload
// (gqlgen graphql-multipart-request-spec format).
func (s *ProductApiTestSuite) executeGraphQLUpload(query string, vars map[string]interface{}, mapData map[string][]string, fileName string, fileContent []byte, authToken string) (*graphtest.GraphQLResponse, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	ops := map[string]interface{}{
		"query":     query,
		"variables": vars,
	}
	opJSON, err := json.Marshal(ops)
	if err != nil {
		return nil, err
	}
	if err := w.WriteField("operations", string(opJSON)); err != nil {
		return nil, err
	}

	mapJSON, err := json.Marshal(mapData)
	if err != nil {
		return nil, err
	}
	if err := w.WriteField("map", string(mapJSON)); err != nil {
		return nil, err
	}

	for key := range mapData {
		part, err := w.CreateFormFile(key, fileName)
		if err != nil {
			return nil, err
		}
		if _, err := part.Write(fileContent); err != nil {
			return nil, err
		}
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	req := httptest.NewRequest(http.MethodPost, "/query", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	rec := httptest.NewRecorder()
	s.handler.ServeHTTP(rec, req)

	var resp graphtest.GraphQLResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *ProductApiTestSuite) TestProductApiLifecycle() {
	// 1. Create
	createQ := `mutation CreateProduct($input: CreateProductInput!) {
  createProduct(input: $input) {
    status
    message
    data {
      id
      name
    }
  }
}`
	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"merchantId":   s.merchantID,
			"categoryId":   s.categoryID,
			"name":         "Test Product",
			"description":  "Test Description",
			"price":        1000,
			"countInStock": 10,
			"brand":        "Test Brand",
			"weight":       1,
			"image":        nil,
		},
	}
	mapData := map[string][]string{"0": {"variables.input.image"}}
	resp, err := s.executeGraphQLUpload(createQ, vars, mapData, "product.jpg", []byte("dummy image content"), s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "createProduct errors: %+v", resp.Errors)
	create, ok := resp.Data["createProduct"].(map[string]interface{})
	s.True(ok, "expected createProduct in data, got %+v", resp.Data)
	createData := create["data"].(map[string]interface{})
	s.Equal("Test Product", createData["name"])
	s.productID = int(createData["id"].(float64))

	// 2. FindById
	findQ := `query {
  findByIdProduct(input: {id: ` + strconv.Itoa(s.productID) + `}) {
    status
    message
    data { id name }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, findQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "findByIdProduct errors: %+v", resp.Errors)
	found, ok := resp.Data["findByIdProduct"].(map[string]interface{})
	s.True(ok, "expected findByIdProduct in data, got %+v", resp.Data)
	foundData := found["data"].(map[string]interface{})
	s.Equal(float64(s.productID), foundData["id"])

	// 3. FindAll
	allQ := `query {
  findAllProduct(input: {page: 1, pageSize: 10}) {
    status
    message
    data { id name }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, allQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "findAllProduct errors: %+v", resp.Errors)

	// 4. FindByActive
	activeQ := `query {
  findByActiveProduct(input: {page: 1, pageSize: 10}) {
    status
    message
    data { id }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, activeQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "findByActiveProduct errors: %+v", resp.Errors)

	// 5. FindByMerchant
	byMerchantQ := `query {
  findByMerchantProduct(input: {
    merchantId: ` + strconv.Itoa(s.merchantID) + `,
    page: 1,
    pageSize: 10,
    minPrice: 0,
    maxPrice: 0
  }) {
    status
    message
    data { id name }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, byMerchantQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "findByMerchantProduct errors: %+v", resp.Errors)

	// 6. FindByCategory
	byCatQ := `query {
  findByCategoryProduct(input: {
    categoryName: "` + s.categoryName + `",
    page: 1,
    pageSize: 10,
    minPrice: 0,
    maxPrice: 0
  }) {
    status
    message
    data { id name }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, byCatQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "findByCategoryProduct errors: %+v", resp.Errors)

	// 7. Update (image optional, use JSON)
	updateQ := `mutation UpdateProduct($input: UpdateProductInput!) {
  updateProduct(input: $input) {
    status
    message
    data { id name }
  }
}`
	updateVars := map[string]interface{}{
		"input": map[string]interface{}{
			"productId":    s.productID,
			"merchantId":   s.merchantID,
			"categoryId":   s.categoryID,
			"name":         "Updated Product",
			"description":  "Updated Description",
			"price":        2000,
			"countInStock": 20,
			"brand":        "Updated Brand",
			"weight":       2,
		},
	}
	resp, err = graphtest.ExecuteGraphQL(s.handler, updateQ, updateVars, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "updateProduct errors: %+v", resp.Errors)
	update, ok := resp.Data["updateProduct"].(map[string]interface{})
	s.True(ok, "expected updateProduct in data, got %+v", resp.Data)
	updateData := update["data"].(map[string]interface{})
	s.Equal("Updated Product", updateData["name"])

	// 8. Trash
	trashQ := `mutation {
  trashedProduct(input: {id: ` + strconv.Itoa(s.productID) + `}) {
    status
    message
    data { id }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, trashQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "trashedProduct errors: %+v", resp.Errors)

	// 9. FindByTrashed
	trashedQ := `query {
  findByTrashedProduct(input: {page: 1, pageSize: 10}) {
    status
    message
    data { id }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, trashedQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "findByTrashedProduct errors: %+v", resp.Errors)

	// 10. Restore
	restoreQ := `mutation {
  restoreProduct(input: {id: ` + strconv.Itoa(s.productID) + `}) {
    status
    message
    data { id }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, restoreQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "restoreProduct errors: %+v", resp.Errors)

	// 11. DeletePermanent
	trashAgainQ := `mutation {
  trashedProduct(input: {id: ` + strconv.Itoa(s.productID) + `}) {
    status
    message
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, trashAgainQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors)

	deleteQ := `mutation {
  deleteProductPermanent(input: {id: ` + strconv.Itoa(s.productID) + `}) {
    status
    message
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, deleteQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "deleteProductPermanent errors: %+v", resp.Errors)

	// 12. RestoreAll
	restoreAllQ := `mutation {
  restoreAllProduct {
    status
    message
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, restoreAllQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "restoreAllProduct errors: %+v", resp.Errors)

	// 13. DeleteAll
	deleteAllQ := `mutation {
  deleteAllProductPermanent {
    status
    message
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, deleteAllQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "deleteAllProductPermanent errors: %+v", resp.Errors)
}

func TestProductApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ProductApiTestSuite))
}
