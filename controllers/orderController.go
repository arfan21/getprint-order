package controllers

import (
	"github.com/arfan21/getprint-order/models"
	"github.com/arfan21/getprint-order/repository"
	"github.com/arfan21/getprint-order/services"
	"github.com/arfan21/getprint-order/utils"
	"github.com/arfan21/getprint-order/validation"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type orderController struct {
	service models.OrderService
}

func NewOrderController(route *echo.Echo, db *gorm.DB) {
	repo := repository.NewOrderRepository(db)
	service := services.NewOrderService(repo)

	ctrl := orderController{service: service}
	route.POST("/order", ctrl.Create)
	route.GET("/order/:id", ctrl.GetById)
	route.GET("/order/user/:id", ctrl.GetByUserId)
	route.PUT("/order/:id", ctrl.Update)
}

func (ctrl *orderController) Create(c echo.Context) error {
	order := new(models.Order)

	if err := c.Bind(order); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.Response("error", err.Error(), nil))
	}

	if err := validation.Validate(*order); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response("error", err, nil))
	}

	err := ctrl.service.Create(order)

	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.Response("error", err.Error(), nil))
	}

	return c.JSON(http.StatusOK, utils.Response("success", nil, order))
}

func (ctrl *orderController) GetById(c echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 32)

	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response("error", err.Error(), nil))
	}

	data, err := ctrl.service.GetByID(uint(id))

	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.Response("error", err.Error(), nil))
	}

	return c.JSON(http.StatusOK, utils.Response("success", nil, data))
}

func (ctrl *orderController) GetByUserId(c echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 32)

	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response("error", err.Error(), nil))
	}

	data, err := ctrl.service.GetByUserID(uint(id))

	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.Response("error", err.Error(), nil))
	}

	return c.JSON(http.StatusOK, utils.Response("success", nil, data))
}

func (ctrl *orderController) Update(c echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 32)

	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response("error", err.Error(), nil))
	}

	order := new(models.Order)
	order.ID = uint(id)
	if err := c.Bind(order); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.Response("error", err.Error(), nil))
	}

	err = ctrl.service.Update(order)

	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.Response("error", err.Error(), nil))
	}

	return c.JSON(http.StatusOK, utils.Response("success", nil, order))
}
