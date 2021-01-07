package repository

import (
	"github.com/arfan21/getprint-order/models"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) models.OrderRepository {
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

	return order, nil
}

func (repo *orderRepository) GetByUserID(userID uint) (*[]models.Order, error) {
	orders := make([]models.Order, 0)
	err := repo.db.Where("user_id = ?", userID).Find(&orders).Error

	if err != nil {
		return nil, err
	}

	return &orders, nil
}

func (repo *orderRepository) Update(data *models.Order) error {
	err := repo.db.Save(&data).Error
	if err != nil {
		return err
	}

	return nil
}