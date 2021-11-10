package basket

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Resource struct {
	domainService Service
}

func ValidateUserFunc(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		token := c.Request().Header.Get("token")

		if len(token) > 0 {
			if token != "ABC" {
				return c.JSON(http.StatusUnauthorized, "You dont have access to api")
			}
			//token len > 0 and token valid
			//user := userService.GetUserInfoByToken()
			//c.Set("currentUserSession",user)
			return next(c)
		}

		return c.JSON(http.StatusUnauthorized, "You dont have access to api")
	}
}
func NewResource(domainService Service) Resource {
	return Resource{domainService: domainService}
}

func RegisterHandlers(instance *echo.Echo, api Resource) {
	instance.GET("/", helloWorld)
	instance.POST("basket", api.createBasket, ValidateUserFunc)
	//                 basket/38640ef4-2033-427d-94e5-718862115d7d
	instance.GET("basket/:id", api.getBasket)
	instance.POST("basket/item", api.addItem)
	// /  basket / 38640ef4-2033-427d-94e5-718862115d7d / item / b0080814-e213-461c-8d09-75f9399320b1
	instance.DELETE("basket/:id/item/:itemid", api.deleteItem)

}

func (receiver *Resource) createBasket(c echo.Context) error {

	request := new(CreateBasketRequest)

	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	result, err := receiver.domainService.Create(c.Request().Context(), request.Buyer)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if result == nil {
		return c.JSON(http.StatusNotFound, "")
	}
	return c.JSON(http.StatusCreated, result)
}

func (receiver *Resource) getBasket(context echo.Context) error {

	id := context.Param("id")
	basket, err := receiver.domainService.Get(context.Request().Context(), id)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}
	if basket == nil {
		return context.JSON(http.StatusNotFound, "")
	}
	// Domain Model Object (such as Basket) ===> BasketDTO
	// Dto => Data Transfer Object
	//dto := &BasketDto{
	//	Id:    basket.Id,
	//	Buyer: basket.BuyerId,
	//}
	return context.JSON(http.StatusOK, basket)
}

func (receiver *Resource) addItem(c echo.Context) error {

	request := new(AddItemRequest)

	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	//if err := c.Validate(request); err != nil {
	//	return c.JSON(http.StatusBadRequest, err.Error())
	//}
	itemId, err := receiver.domainService.AddItem(c.Request().Context(), request.BasketId, request.Sku, request.Quantity, request.Price)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if itemId == "" {
		return c.JSON(http.StatusNotFound, "")
	}
	return c.JSON(http.StatusCreated, itemId)

}

func (receiver *Resource) deleteItem(c echo.Context) error {

	basketid := c.Param("id")
	itemid := c.Param("itemid")

	if len(basketid) == 0 || len(itemid) == 0 {
		return c.JSON(http.StatusBadRequest, "Failed to delete item.")
	}

	if err := receiver.domainService.RemoveItem(c.Request().Context(), basketid, itemid); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusAccepted, "")

}

func helloWorld(context echo.Context) error {

	context.JSON(http.StatusOK, "merhaba dünya")
	return nil
}

//DTO
type CreateBasketRequest struct {
	Buyer string `json:"buyer"`
}
type AddItemRequest struct {
	BasketId string `json:"basketId"`
	Sku      string `json:"sku"`
	Quantity int    `json:"quantity" validate:"gte=1,lt=4"`
	Price    int64  `json:"price"`
}
