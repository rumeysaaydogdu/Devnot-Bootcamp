package basket

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func Test_BasketDomainService(usecase *testing.T) {

	usecase.Run("NewService", func(t *testing.T) {

		tests := []struct {
			name string
			args Repository
			want bool
		}{
			{name: "WithValidArgs_ShouldCreateServiceInstace", args: NewMockRepository(), want: true},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := NewService(tt.args, nil); !tt.want {
					t.Errorf("newService() = %v, want %v", got, tt.want)
				}
			})
		}

	})
	//given => when => then
	// state =>  f(state) => new_state
	usecase.Run("ReadMethods", func(t *testing.T) {
		//Given
		givenBasket := Basket{
			Id:        "ID__1",
			BuyerId:   "Buyer_1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockrepo := &mockRepository{items: []Basket{givenBasket}}
		loadData(mockrepo)
		t.Run("Get Method Tests", func(t *testing.T) {

			type args struct {
				ctx      context.Context
				basketId string
			}
			tests := []struct {
				name       string
				basketId   string
				wantBasket *Basket
				wantErr    bool
			}{
				{name: "WithEmptyBasket_ShouldNotFoundError", basketId: "INVALID_ID", wantBasket: nil, wantErr: false},
				{name: "WithExistBasketID_ShouldReturnBasket", basketId: "ID__1", wantBasket: &givenBasket, wantErr: false},
			}
			receiver := service{
				repo: mockrepo,
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {

					gotBasket, err := receiver.Get(context.TODO(), tt.basketId)
					if (err != nil) != tt.wantErr {
						t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if !reflect.DeepEqual(gotBasket, tt.wantBasket) {
						t.Errorf("Get() gotBasket = %v, want %v", gotBasket, tt.wantBasket)
					}
				})
			}
		})

	})

	usecase.Run("Crud Operations", func(t *testing.T) {

		givenBasket := Basket{
			Id:        "CantDelete",
			BuyerId:   "Buyer",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo := &mockRepository{items: []Basket{givenBasket}}
		loadData(mockRepo)
		domainService := NewService(mockRepo, nil)
		ctx := context.Background()

		t.Run("CreateBasket", func(t *testing.T) {
			t.Run("WithValidBuyer_ShouldBeSuccess", func(t *testing.T) {
				givenBuyer := "Buyer-X"
				got, err := domainService.Create(ctx, givenBuyer)
				assert.Nil(t, err, "You have an errror")
				assert.NotNil(t, got, "the basket is nil")
				assert.Equal(t, got.BuyerId, givenBuyer)
			})
			t.Run("WithBuyerNameIsError_ShouldBeFailed", func(t *testing.T) {
				_, err := domainService.Create(ctx, "error")
				assert.Equal(t, errCRUD, errors.Cause(err))
			})
		})

	})

}
