package services

import (
	"context"
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
	dbx := NewDropbox()
	//File from request body is content/type;base64
	for index, orderDetail := range data.OrderDetails{
		index, orderDetail := index, orderDetail
		errg.Go(func() error{

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

	if err := errg.Wait(); err != nil {
		errDelete := deleteFileDropbox(dbx, data.OrderDetails)

		if errDelete != nil{
			err = errDelete
		}

		return  err
	}
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

func (service *orderService) Update(id uint, data *models.Order) error{
	order, err := service.repo.GetByID(id)
	if err != nil {
		return err
	}

	err = utils.Replace(data, order)
	if err != nil {
		return err
	}

	err = service.repo.Update(data)
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