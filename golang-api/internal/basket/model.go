package basket

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

var (
	ErrNoBuyer  = errors.New("Buyer field can not be null or empty")
	ErrNotFound = errors.New("Item Not Found")
)

const maxQuantityPerProduct = 10

type (
	Basket struct {
		Id        string    `json:"id"`
		BuyerId   string    `json:"buyerId"`
		CreatedAt time.Time `json:"items"`
		UpdatedAt time.Time `json:"updatedAt"`
		Items     []Item    `json:"items"`
	}

	Item struct {
		Id        string `json:"id"`
		Sku       string `json:"sku"`
		UnitPrice int64  `json:"unit_price"`
		Quantity  int    `json:"quantity"`
	}
)

func (receiver *Basket) AddItem(quantity int, price int64, sku string) (*Item, error) {

	if quantity > maxQuantityPerProduct {
		return nil, errors.New(fmt.Sprintf("You can't add more item. Item count can be less then %d", maxQuantityPerProduct))
	}
	index, item := receiver.SearchItemBySku(sku)
	if index > -1 {
		return item, errors.New("Service:Item already added")
	}
	item = &Item{
		Id:        GenerateId(),
		Sku:       sku,
		UnitPrice: price,
		Quantity:  quantity,
	}
	receiver.Items = append(receiver.Items, *item)
	receiver.UpdatedAt = time.Now()

	return item, nil
}
func (receiver *Basket) RemoveItem(itemId string) error {

	if index, _ := receiver.SearchItem(itemId); index != -1 {
		receiver.Items = append(receiver.Items[:index], receiver.Items[index+1:]...)
	} else {
		return ErrNotFound
	}
	return nil
}

func (receiver *Basket) SearchItemBySku(sku string) (int, *Item) {

	for i, n := range receiver.Items {
		if n.Sku == sku {
			return i, &n
		}
	}
	return -1, nil
}
func (receiver *Basket) SearchItem(id string) (int, *Item) {

	for i, n := range receiver.Items {
		if n.Id == id {
			return i, &n
		}
	}
	return -1, nil
}
func Create(buyer string) (*Basket, error) {

	if len(buyer) == 0 {
		return nil, ErrNoBuyer
	}
	return &Basket{
		Id:        GenerateId(),
		BuyerId:   buyer,
		CreatedAt: time.Now(),
		Items:     nil,
	}, nil
}

func GenerateId() string {
	return uuid.New().String()
}
