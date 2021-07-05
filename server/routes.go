package server

import (
	ctrl "github.com/arfan21/getprint-order/app/controllers/http"
	repo "github.com/arfan21/getprint-order/app/repository/mysql"
	"github.com/arfan21/getprint-order/app/services"
	"github.com/arfan21/getprint-order/configs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(mysqlClient configs.Client) *echo.Echo {
	route := echo.New()
	route.Use(middleware.Recover())
	route.Use(middleware.Logger())
	apiV1 := route.Group("/v1")

	// routing order
	orderRepo := repo.NewOrderRepository(mysqlClient)
	orderSrv := services.NewOrderService(orderRepo)
	orderCtrl := ctrl.NewOrderControllers(orderSrv)

	apiOrder := apiV1.Group("/order")
	apiOrder.POST("", orderCtrl.Create)
	apiOrder.GET("/:id", orderCtrl.GetById)
	apiOrder.GET("/user/:id", orderCtrl.GetByUserId)
	apiOrder.GET("/partner/:id", orderCtrl.GetByPartnerId)
	apiOrder.PUT("/:id", orderCtrl.Update)

	// routing cart
	cartRepo := repo.NewCartRepository(mysqlClient)
	cartSrv := services.NewCartServices(mysqlClient, cartRepo)
	cartCtrl := ctrl.NewCartControllers(cartSrv)

	apiCart := apiV1.Group("/cart")
	apiCart.POST("", cartCtrl.Create)
	apiCart.GET("/user/:id", cartCtrl.GetByUserID)
	apiCart.PUT("/updatebatch", cartCtrl.Update)
	apiCart.DELETE("/user/:id", cartCtrl.DeleteByUserID)
	apiCart.DELETE("/:id", cartCtrl.DeleteByID)

	return route
}
