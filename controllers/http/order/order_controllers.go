package order

import (
	"net/http"
	"strconv"

	"github.com/arfan21/getprint-order/models"
	_orderRepo "github.com/arfan21/getprint-order/repository/mysql/order"
	_orderSrv "github.com/arfan21/getprint-order/services/order"
	"github.com/arfan21/getprint-order/utils"
	"github.com/arfan21/getprint-order/validation"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type OrderController interface {
	Routes(route *echo.Echo)
}

type orderController struct {
	service _orderSrv.OrderService
}

func NewOrderController(db *gorm.DB) OrderController {
	repo := _orderRepo.NewOrderRepository(db)
	service := _orderSrv.NewOrderService(repo)

	return &orderController{service: service}
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
	id := c.Param("id")

	data, err := ctrl.service.GetByUserID(id)

	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.Response("error", err.Error(), nil))
	}

	return c.JSON(http.StatusOK, utils.Response("success", nil, data))
}

func (ctrl *orderController) GetByPartnerId(c echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 32)

	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response("error", err.Error(), nil))
	}

	data, err := ctrl.service.GetByPartnerID(uint(id))

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
