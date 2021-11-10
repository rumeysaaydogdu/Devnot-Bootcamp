package basket

import "context"

type Repository interface {
	Create(context.Context, *Basket) error
	Update(context.Context, *Basket) error
	Get(ctx context.Context, id string) (*Basket, error)
	Count(ctx context.Context) (int64, error)
}
