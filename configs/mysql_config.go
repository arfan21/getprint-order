package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/arfan21/getprint-order/app/models"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Client interface{
	Conn() *gorm.DB
	Close() error
}

func NewClient() (Client, error) {
	var DBURL string

	if os.Getenv("DB_PASSWORD") == "" {
		DBURL = fmt.Sprintf("%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	} else {
		DBURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	}

	db, err := gorm.Open(mysql.Open(DBURL), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Order{}, &models.OrderDetail{}, &models.Cart{})
	if err != nil {
		return nil, err
	}
	log.Println("MySql Connected")

	return &client{db}, nil
}

type client struct{
	db *gorm.DB
}

func (c *client) Conn() *gorm.DB{
	return c.db
}

func (c *client) Close() error {
	sqlDB, err := c.db.DB()
	if err == nil{
		return err
	}
	return sqlDB.Close()
}
