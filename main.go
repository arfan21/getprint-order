package main

import (
	"fmt"
	"github.com/arfan21/getprint-order/controllers"
	"github.com/arfan21/getprint-order/utils"
	"github.com/labstack/echo/v4"
	"log"
	"os"
)

func main() {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8000"
	}

	db, err := utils.Connect()
	if err != nil{
		log.Fatal(err)
	}

	route := echo.New()

	controllers.NewOrderController(route, db)

	route.Logger.Fatal(route.Start(fmt.Sprintf(":%s", PORT)))
}
