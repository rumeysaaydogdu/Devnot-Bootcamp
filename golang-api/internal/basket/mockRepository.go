package basket

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"time"
)

var errCRUD = errors.New("Mock: Error crud operation")

type mockRepository struct {
	items []Basket
}

func (m mockRepository) Create(ctx context.Context, basket *Basket) error {
	if basket.BuyerId == "error" {
		return errCRUD
	}
	m.items = append(m.items, *basket)
	return nil
}

func (m mockRepository) Get(ctx context.Context, id string) (*Basket, error) {

	if len(id) == 0 {
		return nil, sql.ErrNoRows
	}
	for _, item := range m.items {
		if item.Id == id {
			return &item, nil
		}
	}
	return nil, nil

}
func (m *mockRepository) Update(ctx context.Context, basket *Basket) error {
	if basket.BuyerId == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.Id == basket.Id {
			m.items[i] = *basket
			break
		}
	}
	return nil
}

func (m mockRepository) Count(ctx context.Context) (int64, error) {

	return int64(len(m.items)), nil
}

func NewMockRepository() Repository {

	return &mockRepository{}
}

func loadData(repo *mockRepository) {

	ctx := context.TODO()
	for i := 0; i < 100; i++ {

		basket := &Basket{
			Id:        fmt.Sprintf("ID_%d", i),
			BuyerId:   fmt.Sprintf("Buyer_ID_%d", i),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		repo.Create(ctx, basket)
	}
}
