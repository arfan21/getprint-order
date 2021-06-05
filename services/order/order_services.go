package order

import (
	"context"
	"reflect"

	"github.com/arfan21/getprint-order/models"
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
	orderRepo   _orderRepo.OrderRepository
	userRepo    _userRepo.UserRepository
	partnerRepo _partnerRepo.PartnerRepository
}

func NewOrderService(repo _orderRepo.OrderRepository) OrderService {
	userRepo := _userRepo.NewUserRepository()
	partnerRepo := _partnerRepo.NewPartnerRepository(context.Background())
	return &orderService{orderRepo: repo, userRepo: userRepo, partnerRepo: partnerRepo}
}

func (service *orderService) Create(data *models.Order) error {
	errg, ctx := errgroup.WithContext(context.Background())

	partnerChan := make(chan *_partnerRepo.PartnerResponse)
	defer close(partnerChan)

	//find user
	errg.Go(func() error {
		_, err := service.userRepo.GetUserByID(ctx, data.UserID.String())
		if err != nil {
			return err
		}

		return nil
	})

	//find partner
	errg.Go(func() error {
		data, err := service.partnerRepo.GetPartnerByID(data.PartnerID)
		if err != nil {
			partnerChan <- nil

			return err
		}
		partnerChan <- data

		return nil
	})

	partner := <-partnerChan

	if err := errg.Wait(); err != nil {
		return err
	}

	totalPrice := countPrice(data.OrderDetails, partner)
	data.TotalPrice = totalPrice

	err := service.orderRepo.Create(data)
	if err != nil {

		return err
	}

	return nil
}

func (service *orderService) GetByID(id uint) (*models.Order, error) {
	order, err := service.orderRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (service *orderService) GetByUserID(userID string) (*[]models.Order, error) {
	userIDUUID, err := uuid.FromString(userID)
	if err != nil {
		return nil, err
	}

	orders, err := service.orderRepo.GetByUserID(userIDUUID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (service *orderService) GetByPartnerID(partnerID uint) (*[]models.Order, error) {
	orders, err := service.orderRepo.GetByPartnerID(partnerID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (service *orderService) Update(data *models.Order) error {
	order, err := service.orderRepo.GetByID(data.ID)
	if err != nil {
		return err
	}

	err = service.orderRepo.Update(data)
	if err != nil {
		return err
	}

	return utils.Replace(*order, data)
}

func countPrice(data []models.OrderDetail, partnerResponse *_partnerRepo.PartnerResponse) uint {
	var totalPrice uint = 0

	dataPartner := partnerResponse.Data
	elements := reflect.ValueOf(&dataPartner).Elem()

	for _, order := range data {
		for i := 0; i < elements.NumField(); i++ {
			price := elements.Field(i)
			nameJson := elements.Type().Field(i).Tag.Get("json")
			if order.OrderType == nameJson {
				totalPrice += (order.Qty * uint(price.Int()))
			}
		}
	}

	return totalPrice
}
