package basket

import (
	"testing"
)

func TestBasketModel(usecase *testing.T) {

	usecase.Run("CreateBasket", func(t *testing.T) {

		tests := []struct {
			name    string
			buyer   string
			wantErr bool
		}{
			{"WithBasketHasBuyer_ShouldSuccess", "buyer", false},
			{"WithBasketBuyerIsEmpty_ShouldFailed", "", true},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := Create(tt.buyer)
				if (err != nil) != tt.wantErr {
					t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			})
		}

	})
	usecase.Run("AddItem", func(t *testing.T) {

	})

}
