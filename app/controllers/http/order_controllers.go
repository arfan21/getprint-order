package http

import (
	"net/http"
	"strconv"

	"github.com/arfan21/getprint-order/app/helpers"
	models2 "github.com/arfan21/getprint-order/app/models"
	"github.com/arfan21/getprint-order/app/services"
	"github.com/arfan21/getprint-order/validation"
	"github.com/labstack/echo/v4"
)

type OrderControllers interface {
	Create(c echo.Context) error
	GetById(c echo.Context) error
	GetByUserId(c echo.Context) error
	GetByPartnerId(c echo.Context) error
	Update(c echo.Context) error
}

type orderControllers struct {
	service services.OrderService
}

func NewOrderControllers(service services.OrderService) OrderControllers {
	return &orderControllers{service: service}
}

func (ctrl *orderControllers) Create(c echo.Context) error {
	order := new(models2.Order)

	if err := c.Bind(order); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, helpers.Response("error", err.Error(), nil))
	}

	if err := validation.Validate(*order); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.Response("error", err, nil))
	}

	err := ctrl.service.Create(order)

	if err != nil {
		return c.JSON(helpers.GetStatusCode(err), helpers.Response("error", err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helpers.Response("success", nil, order))
}

func (ctrl *orderControllers) GetById(c echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 32)

	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.Response("error", err.Error(), nil))
	}

	data, err := ctrl.service.GetByID(uint(id))

	if err != nil {
		return c.JSON(helpers.GetStatusCode(err), helpers.Response("error", err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helpers.Response("success", nil, data))
}

func (ctrl *orderControllers) GetByUserId(c echo.Context) error {
	id := c.Param("id")

	data, err := ctrl.service.GetByUserID(id)

	if err != nil {
		return c.JSON(helpers.GetStatusCode(err), helpers.Response("error", err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helpers.Response("success", nil, data))
}

func (ctrl *orderControllers) GetByPartnerId(c echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 32)

	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.Response("error", err.Error(), nil))
	}

	data, err := ctrl.service.GetByPartnerID(uint(id))

	if err != nil {
		return c.JSON(helpers.GetStatusCode(err), helpers.Response("error", err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helpers.Response("success", nil, data))
}

func (ctrl *orderControllers) Update(c echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 32)

	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.Response("error", err.Error(), nil))
	}

	order := new(models2.Order)
	order.ID = uint(id)
	if err := c.Bind(order); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, helpers.Response("error", err.Error(), nil))
	}

	err = ctrl.service.Update(order)

	if err != nil {
		return c.JSON(helpers.GetStatusCode(err), helpers.Response("error", err.Error(), nil))
	}

	return c.JSON(http.StatusOK, helpers.Response("success", nil, order))
}
