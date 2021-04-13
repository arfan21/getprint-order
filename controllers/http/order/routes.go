package order

import "github.com/labstack/echo/v4"

func (ctrl orderController) Routes(route *echo.Echo) {
	route.POST("/order", ctrl.Create)
	route.GET("/order/:id", ctrl.GetById)
	route.GET("/order/user/:id", ctrl.GetByUserId)
	route.GET("/order/partner/:id", ctrl.GetByPartnerId)
	route.PUT("/order/:id", ctrl.Update)
}
