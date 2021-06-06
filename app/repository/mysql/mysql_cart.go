package mysql

import (
	models2 "github.com/arfan21/getprint-order/app/models"
	"github.com/arfan21/getprint-order/config/mysql"
	uuid "github.com/satori/go.uuid"
)

type CartRepository interface {
	Create(data *models2.Cart) error
	GetByUserID(userID uuid.UUID) (*[]models2.Cart, error)
	GetByID(id uint) (*models2.Cart, error)
	Update(data *models2.Cart) error
	DeleteByCartID(id uint) error
	DeleteByUserID(userID uuid.UUID) error
}

type cartRepository struct {
	db mysql.Client
}

func NewCartRepository(db mysql.Client) CartRepository {
	return &cartRepository{db}
}

func (repo *cartRepository) Create(data *models2.Cart) error {
	return repo.db.Conn().Create(data).Error
}

func (repo *cartRepository) GetByUserID(userID uuid.UUID) (*[]models2.Cart, error) {
	carts := make([]models2.Cart, 0)
	err := repo.db.Conn().Where("user_id = ?", userID).Find(&carts).Error
	if err != nil {
		return nil, err
	}
	return &carts, nil
}

func (repo *cartRepository) GetByID(id uint) (*models2.Cart, error) {
	cart := new(models2.Cart)
	err := repo.db.Conn().First(cart, id).Error
	if err != nil{
		return nil, err
	}
	return cart, nil
}

func (repo *cartRepository) Update(data *models2.Cart) error {
	return repo.db.Conn().Save(data).Error
}

func (repo *cartRepository) DeleteByCartID(id uint) error {
	return repo.db.Conn().Unscoped().Delete(&models2.Cart{}, id).Error
}

func (repo *cartRepository) DeleteByUserID(userID uuid.UUID) error {
	return repo.db.Conn().Unscoped().Where("user_id = ?",userID).Delete(&models2.Cart{}).Error
}
