package basket

import (
	"context"
	"fmt"
	rabbit "github.com/alperhankendi/golang-api/pkg/rabbitmq"
	"github.com/pkg/errors"
)

// Service encapsulates use case logic for the basket.
type Service interface {
	Create(ctx context.Context, buyer string) (*Basket, error)
	Get(ctx context.Context, basketId string) (*Basket, error)

	AddItem(context.Context, string, string, int, int64) (string, error)
	RemoveItem(context.Context, string, string) error
}
type service struct {
	repo         Repository //interface, not mongo collection or mongo struct
	rabbitClient *rabbit.Client
}

/// Application (API).Create() ===> DomainService.Create() ==> DB.save

//Create Basket Method
func (receiver service) Create(ctx context.Context, buyer string) (*Basket, error) {

	//business logic
	basket, err := Create(buyer)
	if err != nil {
		return nil, errors.Wrap(err, "Service:Failed to create basket object")
	}
	err = receiver.repo.Create(ctx, basket)
	if err != nil {
		return nil, errors.Wrap(err, "Service:Failed to create basket")
	}

	if err = receiver.rabbitClient.Publish(context.TODO(), "", *basket); err != nil {
		fmt.Println("error:", err)
	}

	return basket, nil
}

//ServiceFactory
func NewService(r Repository, rclient *rabbit.Client) Service {

	if r == nil {
		return nil
	}
	return &service{
		repo:         r,
		rabbitClient: rclient,
	}
}

func (receiver service) AddItem(ctx context.Context, basketId string, sku string, qty int, price int64) (string, error) {

	basket, err := receiver.repo.Get(ctx, basketId)
	if err != nil {
		return "", errors.Wrap(err, "Service:failed to retrieve basket")
	}
	if basket == nil {
		return "", errors.New("Service:Basket not found")
	}

	//model.go AddItem // domain model
	item, err := basket.AddItem(qty, price, sku)
	if err != nil {
		return "", errors.Wrap(err, "failed to add item to basket.")
	}
	//update basket
	if err = receiver.repo.Update(ctx, basket); err != nil {
		return "", errors.Wrap(err, "failed to update basket")
	}

	return item.Id, nil
}

func (receiver service) RemoveItem(ctx context.Context, basketId string, itemId string) error {

	basket, err := receiver.repo.Get(ctx, basketId)
	if err != nil {
		return errors.Wrap(err, "Service:failed to retrieve basket")
	}
	if basket == nil {
		return errors.New("Service:Basket not found")
	}
	if err = basket.RemoveItem(itemId); err != nil {
		return errors.Wrap(err, "failed to add item to basket.")
	}
	if err = receiver.repo.Update(ctx, basket); err != nil {
		return errors.Wrap(err, "failed to delete item")
	}
	return nil
}

func (receiver service) Get(ctx context.Context, basketId string) (basket *Basket, err error) {

	basket, err = receiver.repo.Get(ctx, basketId)
	if err != nil {
		err = errors.Wrapf(err, "Failed to retrieve basket data. Basket Id: %s", basketId)
	}
	return
}
