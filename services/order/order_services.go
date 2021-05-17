package order

import (
	"context"
	"fmt"
	"strings"

	"github.com/arfan21/getprint-order/models"
	_dropboxRepo "github.com/arfan21/getprint-order/repository/dropbox"
	_orderRepo "github.com/arfan21/getprint-order/repository/mysql/order"
	_partnerRepo "github.com/arfan21/getprint-order/repository/partner"
	_userRepo "github.com/arfan21/getprint-order/repository/user"
	"github.com/arfan21/getprint-order/utils"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/sync/errgroup"
)

type OrderService interface {
	Create(data *models.Order) error
	GetByID(id uint) (*models.Order, error)
	GetByUserID(userID string) (*[]models.Order, error)
	GetByPartnerID(partnerID uint) (*[]models.Order, error)
	Update(data *models.Order) error
}

type orderService struct {
	repo _orderRepo.OrderRepository
}

func NewOrderService(repo _orderRepo.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (service *orderService) Create(data *models.Order) error {
	errg, ctx := errgroup.WithContext(context.Background())

	userRepo := _userRepo.NewUserRepository(ctx)
	partnerRepo := _partnerRepo.NewPartnerRepository(ctx)

	partnerChan := make(chan *_partnerRepo.PartnerResponse)
	defer close(partnerChan)

	//find user
	errg.Go(func() error {
		fmt.Println("check user")
		_, err := userRepo.GetUserByID(data.UserID.String())
		fmt.Println("error  user :", err)
		if err != nil {
			return err
		}

		return nil
	})

	//find partner
	errg.Go(func() error {
		data, err := partnerRepo.GetPartnerByID(data.PartnerID)
		fmt.Println("mengirim data partner")
		defer fmt.Println("selesai mengirim data partner")
		if err != nil {
			partnerChan <- nil

			return err
		}
		partnerChan <- data

		return nil
	})
	partner := <-partnerChan
	fmt.Println("menerima data")
	fmt.Println(partner)
	if err := errg.Wait(); err != nil {
		return err
	}

	errg2, ctx2 := errgroup.WithContext(context.Background())
	dropboxRepo := _dropboxRepo.NewDropboxRepository(ctx2)
	//Uploading file and file from request body is content/type;base64
	for index, orderDetail := range data.OrderDetails {
		index, orderDetail := index, orderDetail
		errg2.Go(func() error {
			buffer, filename, err := utils.GetFileBufferAndFileName(orderDetail.File)
			if err != nil {
				return models.ErrUnprocessableEntity
			}

			path, err := dropboxRepo.Uploader(filename, buffer)
			if err != nil {
				return err
			}

			data.OrderDetails[index].Path = path

			sharedLink, err := dropboxRepo.CreateSharedLink(path)
			if err != nil {
				return err
			}
			data.OrderDetails[index].File = sharedLink

			return nil
		})
	}

	if err := errg2.Wait(); err != nil {
		errDelete := deleteFileDropbox(dropboxRepo, data.OrderDetails)

		if errDelete != nil {
			err = errDelete
		}

		return err
	}

	totalPrice := countPrice(data.OrderDetails, partner)
	data.TotalPrice = totalPrice

	err := service.repo.Create(data)
	if err != nil {
		errDelete := deleteFileDropbox(dropboxRepo, data.OrderDetails)

		if errDelete != nil {
			err = errDelete
		}

		return err
	}

	return nil
}

func (service *orderService) GetByID(id uint) (*models.Order, error) {
	order, err := service.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (service *orderService) GetByUserID(userID string) (*[]models.Order, error) {
	userIDUUID , err := uuid.FromString(userID)
	if err != nil {
		return nil, err
	}

	orders, err := service.repo.GetByUserID(userIDUUID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (service *orderService) GetByPartnerID(partnerID uint) (*[]models.Order, error) {
	orders, err := service.repo.GetByPartnerID(partnerID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (service *orderService) Update(data *models.Order) error {
	order, err := service.repo.GetByID(data.ID)
	if err != nil {
		return err
	}

	err = service.repo.Update(data)
	if err != nil {
		return err
	}

	return utils.Replace(*order, data)
}

func deleteFileDropbox(dbx _dropboxRepo.Dropbox, data []models.OrderDetail) error {
	errg, _ := errgroup.WithContext(context.Background())

	for _, orderDetail := range data {
		orderDetail := orderDetail
		if !strings.Contains(orderDetail.File, "base64") {
			errg.Go(func() error {
				err := dbx.Delete(orderDetail.Path)
				if err != nil {
					return err
				}
				return nil
			})
		}
	}

	return errg.Wait()
}

func countPrice(data []models.OrderDetail, res *_partnerRepo.PartnerResponse) uint {
	printPrice := uint(res.Data.Print.Int64)
	scanPrice := uint(res.Data.Scan.Int64)
	photocopyPrice := uint(res.Data.Fotocopy.Int64)

	var totalPrice uint = 0

	for _, order := range data {
		if order.PrintQty != 0 {
			totalPrice = totalPrice + (order.PrintQty * printPrice)
		}

		if order.ScanQty != 0 {
			totalPrice = totalPrice + (order.ScanQty * scanPrice)
		}

		if order.PhotocopyQty != 0 {
			totalPrice = totalPrice + (order.PhotocopyQty * photocopyPrice)
		}

	}

	return totalPrice
}
