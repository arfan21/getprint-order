package main

import (
	"fmt"
	"log"
	"os"

	_orderCtrl "github.com/arfan21/getprint-order/controllers/http/order"
	"github.com/arfan21/getprint-order/utils"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

func main() {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8000"
	}

	db, err := utils.Connect()
	if err != nil {
		log.Fatal(err)
	}

	route := echo.New()

	orderCtrl := _orderCtrl.NewOrderController(db)
	orderCtrl.Routes(route)

	route.Logger.Fatal(route.Start(fmt.Sprintf(":%s", PORT)))
}
