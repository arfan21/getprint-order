package mysql

import (
	"context"

	models2 "github.com/arfan21/getprint-order/app/models"
	"github.com/arfan21/getprint-order/configs"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/sync/errgroup"
)

type OrderRepository interface {
	Create(data *models2.Order) error
	GetByID(id uint) (*models2.Order, error)
	GetByUserID(userID uuid.UUID) (*[]models2.Order, error)
	GetByPartnerID(partnerID uint) (*[]models2.Order, error)
	Update(data *models2.Order) error
}

type orderRepository struct {
	db configs.Client
}

func NewOrderRepository(db configs.Client) OrderRepository {
	return &orderRepository{db: db}
}

func (repo *orderRepository) Create(data *models2.Order) error {
	err := repo.db.Conn().Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *orderRepository) GetByID(id uint) (*models2.Order, error) {
	order := new(models2.Order)
	err := repo.db.Conn().First(&order, id).Error
	if err != nil {
		return nil, err
	}

	orderDetails := make([]models2.OrderDetail, 0)
	err = repo.db.Conn().Where("order_id=?", id).Find(&orderDetails).Error
	if err != nil {
		return nil, err
	}
	order.OrderDetails = orderDetails

	return order, nil
}

func (repo *orderRepository) GetByUserID(userID uuid.UUID) (*[]models2.Order, error) {
	orders := make([]models2.Order, 0)
	err := repo.db.Conn().Debug().Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	errg, ctx := errgroup.WithContext(context.Background())
	for index, order := range orders {
		index, order := index, order
		errg.Go(func() error {
			orderDetails := make([]models2.OrderDetail, 0)
			err := repo.db.Conn().WithContext(ctx).Where("order_id=?", order.ID).Find(&orderDetails).Error

			if err != nil {
				return err
			}
			orders[index].OrderDetails = orderDetails
			return nil
		})
	}

	if err := errg.Wait(); err != nil {
		return nil, err
	}

	return &orders, nil
}

func (repo *orderRepository) GetByPartnerID(partnerID uint) (*[]models2.Order, error) {
	orders := make([]models2.Order, 0)
	err := repo.db.Conn().Where("partner_id = ?", partnerID).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	errg, ctx := errgroup.WithContext(context.Background())
	for index, order := range orders {
		index, order := index, order
		errg.Go(func() error {
			orderDetails := make([]models2.OrderDetail, 0)
			err := repo.db.Conn().WithContext(ctx).Where("order_id=?", order.ID).Find(&orderDetails).Error

			if err != nil {
				return err
			}
			orders[index].OrderDetails = orderDetails
			return nil
		})
	}

	if err := errg.Wait(); err != nil {
		return nil, err
	}

	return &orders, nil
}

func (repo *orderRepository) Update(data *models2.Order) error {
	err := repo.db.Conn().Model(&data).Update("status", data.Status).Error
	if err != nil {
		return err
	}

	return nil
}
