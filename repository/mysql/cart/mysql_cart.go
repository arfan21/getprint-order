package cart

import (
	"github.com/arfan21/getprint-order/models"
	"gorm.io/gorm"
)

type CartRepository interface{}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db}
}

func (repo *cartRepository) Create(data *models.Cart) error {
	return repo.db.Create(data).Error
}
