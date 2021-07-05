package mysql

import (
	"github.com/arfan21/getprint-order/app/models"
	"github.com/arfan21/getprint-order/configs"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type CartRepository interface {
	Create(data []models.Cart) error
	GetByUserID(userID uuid.UUID) (*[]models.Cart, error)
	GetByID(id uint) (*models.Cart, error)
	UpdateByID(data *models.Cart) error
	UpdateByIDwithTx(tx *gorm.DB, data models.Cart) error
	DeleteByCartID(id uint) error
	DeleteByUserID(userID uuid.UUID) error
}

type cartRepository struct {
	db configs.Client
}

func NewCartRepository(db configs.Client) CartRepository {
	return &cartRepository{db}
}

func (repo *cartRepository) Create(data []models.Cart) error {
	return repo.db.Conn().Debug().Create(data).Error
}

func (repo *cartRepository) GetByUserID(userID uuid.UUID) (*[]models.Cart, error) {
	carts := make([]models.Cart, 0)
	err := repo.db.Conn().Where("user_id = ?", userID).Find(&carts).Error
	if err != nil {
		return nil, err
	}
	return &carts, nil
}

func (repo *cartRepository) GetByID(id uint) (*models.Cart, error) {
	cart := new(models.Cart)
	err := repo.db.Conn().First(cart, id).Error
	if err != nil {
		return nil, err
	}
	return cart, nil
}

func (repo *cartRepository) UpdateByID(data *models.Cart) error {
	return repo.db.Conn().Save(data).Error
}

func (repo *cartRepository) UpdateByIDwithTx(tx *gorm.DB, data models.Cart) error {
	return repo.db.Conn().Save(data).Error
}

func (repo *cartRepository) DeleteByCartID(id uint) error {
	return repo.db.Conn().Unscoped().Delete(&models.Cart{}, id).Error
}

func (repo *cartRepository) DeleteByUserID(userID uuid.UUID) error {
	return repo.db.Conn().Unscoped().Where("user_id = ?", userID).Delete(&models.Cart{}).Error
}
