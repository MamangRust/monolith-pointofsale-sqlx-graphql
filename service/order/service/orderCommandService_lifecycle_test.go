package service

import (
	"context"
	"errors"
	"testing"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
)

type lifecycleStockCall struct {
	productID int
	delta     int
}

type lifecycleStockRepository struct {
	calls       []lifecycleStockCall
	failProduct int
	failOn      int
}

func (r *lifecycleStockRepository) AdjustProductStock(_ context.Context, productID, delta int) (*db.Product, error) {
	r.calls = append(r.calls, lifecycleStockCall{productID: productID, delta: delta})
	if !r.shouldFail(productID) {
		return &db.Product{ProductID: int32(productID)}, nil
	}
	return nil, errors.New("adjust stock failed")
}

func (r *lifecycleStockRepository) DecrementProductCountStock(ctx context.Context, productID, quantity int) (*db.Product, error) {
	return r.AdjustProductStock(ctx, productID, -quantity)
}

func (r *lifecycleStockRepository) IncrementProductCountStock(ctx context.Context, productID, quantity int) (*db.Product, error) {
	return r.AdjustProductStock(ctx, productID, quantity)
}

func (r *lifecycleStockRepository) shouldFail(productID int) bool {
	if r.failProduct != productID {
		return false
	}
	return r.failOn == 0 || len(r.calls) == r.failOn
}

// lifecycleProductQueryRepository backs the post-failure stock re-read used to
// label insufficient-stock failures.
type lifecycleProductQueryRepository struct{}

func (r *lifecycleProductQueryRepository) FindById(_ context.Context, id int) (*db.Product, error) {
	return &db.Product{ProductID: int32(id), CountInStock: 0}, nil
}

func TestIncrementStockForOrderItemsCompensatesPartialMutation(t *testing.T) {
	repo := &lifecycleStockRepository{failProduct: 2}
	service := &orderCommandService{productCommandRepository: repo}
	items := []*db.OrderItem{
		{ProductID: 1, Quantity: 3},
		{ProductID: 2, Quantity: 4},
	}

	_, err := service.incrementStockForOrderItems(context.Background(), items)
	if err == nil {
		t.Fatal("expected increment failure")
	}

	want := []lifecycleStockCall{
		{productID: 1, delta: 3},
		{productID: 2, delta: 4},
		{productID: 1, delta: -3},
	}
	assertLifecycleStockCalls(t, want, repo.calls)
}

func TestDecrementStockForOrderItemsCompensatesPartialMutation(t *testing.T) {
	repo := &lifecycleStockRepository{failProduct: 2}
	service := &orderCommandService{
		productCommandRepository: repo,
		productQueryRepository:   &lifecycleProductQueryRepository{},
	}
	items := []*db.OrderItem{
		{ProductID: 1, Quantity: 2},
		{ProductID: 2, Quantity: 5},
	}

	_, err := service.decrementStockForOrderItems(context.Background(), items)
	if err == nil {
		t.Fatal("expected decrement failure")
	}

	want := []lifecycleStockCall{
		{productID: 1, delta: -2},
		{productID: 2, delta: -5},
		{productID: 1, delta: 2},
	}
	assertLifecycleStockCalls(t, want, repo.calls)
}

func assertLifecycleStockCalls(t *testing.T, want, got []lifecycleStockCall) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("unexpected call count: got %d, want %d; calls=%v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("unexpected call at index %d: got %+v, want %+v", i, got[i], want[i])
		}
	}
}
