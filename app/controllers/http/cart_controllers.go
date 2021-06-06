package http

import (
	"net/http"
	"strconv"

	"github.com/arfan21/getprint-order/app/helpers"
	"github.com/arfan21/getprint-order/app/models"
	"github.com/arfan21/getprint-order/app/services"
	"github.com/labstack/echo/v4"
)

type CartControllers interface{
	Create(c echo.Context) error
	GetByUserID(c echo.Context) error
	Update(c echo.Context) error
	DeleteByID(c echo.Context) error
	DeleteByUserID(c echo.Context) error
}

type cartControllers struct{
	cartSrv services.CartServices
}

func NewCartControllers(cartSrv services.CartServices) CartControllers{
	return &cartControllers{cartSrv}
}

func (ctrl *cartControllers) Create(c echo.Context) error {
	cart := new(models.Cart)

	if err := c.Bind(cart); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, helpers.Response("error", err.Error(), nil))
	}

	err := ctrl.cartSrv.Create(cart)
	if err != nil {
		return c.JSON(helpers.GetStatusCode(err), helpers.Response("error", err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helpers.Response("success", nil, cart))
}

func (ctrl *cartControllers) GetByUserID(c echo.Context) error {
	id := c.Param("id")

	data, err := ctrl.cartSrv.GetByUserID(id)
	if err != nil {
		return c.JSON(helpers.GetStatusCode(err), helpers.Response("error", err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helpers.Response("success", nil, data))
}

func (ctrl *cartControllers) Update(c echo.Context) error {
	cart := new(models.Cart)

	if err := c.Bind(cart); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, helpers.Response("error", err.Error(), nil))
	}

	err := ctrl.cartSrv.Create(cart)
	if err != nil {
		return c.JSON(helpers.GetStatusCode(err), helpers.Response("error", err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helpers.Response("success", nil, cart))
}

func (ctrl *cartControllers) DeleteByID(c echo.Context) error{
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.Response("error", err.Error(), nil))
	}

	err = ctrl.cartSrv.DeleteByID(uint(id))
	if err != nil {
		return c.JSON(helpers.GetStatusCode(err), helpers.Response("error", err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helpers.Response("success", nil, nil))
}

func (ctrl *cartControllers) DeleteByUserID(c echo.Context) error{
	id := c.Param("id")

	err := ctrl.cartSrv.DeleteByUserID(id)
	if err != nil {
		return c.JSON(helpers.GetStatusCode(err), helpers.Response("error", err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helpers.Response("success", nil, nil))
}