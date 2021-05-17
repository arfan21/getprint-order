package order

import (
	"context"

	"github.com/arfan21/getprint-order/models"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(data *models.Order) error
	GetByID(id uint) (*models.Order, error)
	GetByUserID(userID uuid.UUID) (*[]models.Order, error)
	GetByPartnerID(partnerID uint) (*[]models.Order, error)
	Update(data *models.Order) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (repo *orderRepository) Create(data *models.Order) error {
	err := repo.db.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *orderRepository) GetByID(id uint) (*models.Order, error) {
	order := new(models.Order)
	err := repo.db.First(&order, id).Error
	if err != nil {
		return nil, err
	}

	orderDetails := make([]models.OrderDetail, 0)
	err = repo.db.Where("order_id=?", id).Find(&orderDetails).Error
	if err != nil {
		return nil, err
	}
	order.OrderDetails = orderDetails

	return order, nil
}

func (repo *orderRepository) GetByUserID(userID uuid.UUID) (*[]models.Order, error) {
	orders := make([]models.Order, 0)
	err := repo.db.Debug().Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	errg, ctx := errgroup.WithContext(context.Background())
	for index, order := range orders {
		index, order := index, order
		errg.Go(func() error {
			orderDetails := make([]models.OrderDetail, 0)
			err := repo.db.WithContext(ctx).Where("order_id=?", order.ID).Find(&orderDetails).Error

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

func (repo *orderRepository) GetByPartnerID(partnerID uint) (*[]models.Order, error) {
	orders := make([]models.Order, 0)
	err := repo.db.Where("partner_id = ?", partnerID).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	errg, ctx := errgroup.WithContext(context.Background())
	for index, order := range orders {
		index, order := index, order
		errg.Go(func() error {
			orderDetails := make([]models.OrderDetail, 0)
			err := repo.db.WithContext(ctx).Where("order_id=?", order.ID).Find(&orderDetails).Error

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

func (repo *orderRepository) Update(data *models.Order) error {
	err := repo.db.Model(&data).Update("status", data.Status).Error
	if err != nil {
		return err
	}

	return nil
}
