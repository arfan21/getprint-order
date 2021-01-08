package services

import (
	"context"
	"fmt"
	"github.com/arfan21/getprint-order/models"
	"github.com/arfan21/getprint-order/utils"
	"golang.org/x/sync/errgroup"
	"strings"
)

type orderService struct {
	repo models.OrderRepository
}

func NewOrderService(repo models.OrderRepository) models.OrderService{
	return &orderService{repo: repo}
}

func (service *orderService) Create(data *models.Order) error{
	errg, _ := errgroup.WithContext(context.Background())

	partnerChan := make(chan map[string]interface{})

	//find user
	errg.Go(func() error{
		_, err := GetUser(data.UserID)

		if err != nil {
			return err
		}

		return nil
	})

	//find partner
	errg.Go(func() error{
		data, err := GetPartner(data.PartnerID)

		if err != nil {
			close(partnerChan)
			return err
		}
		partnerChan <- data
		close(partnerChan)
		return nil
	})
	partner := <-partnerChan

	if err := errg.Wait(); err != nil {
		return err
	}

	errg2, _ := errgroup.WithContext(context.Background())

	dbx := NewDropbox()
	//Uploading file and file from request body is content/type;base64
	for index, orderDetail := range data.OrderDetails{
		index, orderDetail := index, orderDetail
		errg2.Go(func() error{

			fmt.Println("Upload ", index)
			buffer, filename, err := utils.GetFileBufferAndFileName(orderDetail.File)
			if err != nil{
				return models.ErrUnprocessableEntity
			}

			res, err := dbx.Upload(filename, buffer)
			if err != nil {
				return err
			}

			path := res["path_lower"].(string)
			data.OrderDetails[index].Path = path

			sharedLink, err := dbx.CreateSharedLink(path)
			if err != nil {
				return err
			}
			data.OrderDetails[index].File = sharedLink

			return nil
		})
	}

	if err := errg2.Wait(); err != nil {
		errDelete := deleteFileDropbox(dbx, data.OrderDetails)

		if errDelete != nil{
			err = errDelete
		}

		return  err
	}
	price := partner["data"].(map[string]interface{})["price"].(map[string]interface{})
	totalPrice := countPrice(data.OrderDetails, price)
	data.TotalPrice = totalPrice

	err := service.repo.Create(data)
	if err != nil {
		errDelete := deleteFileDropbox(dbx, data.OrderDetails)

		if errDelete != nil{
			err = errDelete
		}

		return err
	}

	return nil
}

func (service *orderService) GetByID(id uint) (*models.Order, error){
	order, err := service.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (service *orderService) GetByUserID(userID uint) (*[]models.Order, error){
	orders, err := service.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (service *orderService) Update(data *models.Order) error{
	order, err := service.repo.GetByID(data.ID)
	if err != nil {
		return err
	}

	err = service.repo.Update(data)
	if err != nil {
		return err
	}

	err = utils.Replace(*order, data)
	if err != nil {
		return err
	}

	return nil
}

func deleteFileDropbox(dbx Dropbox, data []models.OrderDetail) error{
	errg, _ := errgroup.WithContext(context.Background())

	for _, orderDetail := range data{
		orderDetail := orderDetail
		if !strings.Contains(orderDetail.File, "base64"){
			errg.Go(func()error{
				err := dbx.Delete(orderDetail.Path)
				if err != nil{
					return err
				}
				return nil
			})
		}
	}

	if err := errg.Wait(); err != nil {
		return  err
	}

	return nil
}

func countPrice(data []models.OrderDetail, price map[string]interface{}) uint {
	printPrice := uint(price["print"].(float64))
	scanPrice := uint(price["print"].(float64))
	photocopyPrice := uint(price["print"].(float64))

	var totalPrice uint = 0

	for _, order := range data{
		if order.PrintQty != 0 {
			totalPrice = totalPrice + (order.PrintQty * printPrice)
		}

		if order.ScanQty != 0{
			totalPrice = totalPrice + (order.ScanQty * scanPrice)
		}

		if order.PhotocopyQty != 0{
			totalPrice = totalPrice + (order.PhotocopyQty * photocopyPrice)
		}

	}

	return totalPrice
}