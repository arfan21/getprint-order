package controllers

import (
	"github.com/arfan21/getprint-order/models"
	"github.com/arfan21/getprint-order/repository"
	"github.com/arfan21/getprint-order/services"
	"github.com/arfan21/getprint-order/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

type orderController struct {
	service models.OrderService
}

func NewOrderController(route *echo.Echo, db *gorm.DB) {
	repo := repository.NewOrderRepository(db)
	service := services.NewOrderService(repo)

	ctrl := orderController{service: service}
	route.POST("/order", ctrl.Create)
}

func (ctrl *orderController) Create(c echo.Context) error {
	order := new(models.Order)

	if err := c.Bind(order); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.Response("error", err.Error(), nil))
	}

	if err := order.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response("error", err, nil))
	}

	err := ctrl.service.Create(order)

	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.Response("error", err.Error(), nil))
	}

	return c.JSON(http.StatusOK, utils.Response("success", "success create order", order))
}

func (ctrl *orderController) GetById(c echo.Context) error {
	return nil
}

func (ctrl *orderController) GetByUserId(c echo.Context) error {
	return nil
}

func (ctrl *orderController) Update(c echo.Context) error {
	return nil
}
