package services

import (
	"context"

	"github.com/arfan21/getprint-order/app/models"
	_mediaRepo "github.com/arfan21/getprint-order/app/repository/media"
	repo "github.com/arfan21/getprint-order/app/repository/mysql"
	uuid "github.com/satori/go.uuid"
)

type CartServices interface{
	Create(data *models.Cart) error
	GetByUserID(userid string) (*[]models.Cart, error)
	Update(data *models.Cart) error
	DeleteByID(id uint) error
	DeleteByUserID(userid string) error
}

type cartServices struct {
	cartRepo  repo.CartRepository
	mediaRepo _mediaRepo.MediaRepository
}

func NewCartServices(cartRepo repo.CartRepository) CartServices {
	mediaRepo := _mediaRepo.NewMediaRepository()
	return &cartServices{cartRepo,mediaRepo}
}

func (srv *cartServices) Create(data *models.Cart) error {
	return srv.cartRepo.Create(data)
}

func (srv *cartServices) GetByUserID(userid string) (*[]models.Cart, error){
	userIdUUID, err:= uuid.FromString(userid)
	if err != nil {
		return nil, err
	}
	return srv.cartRepo.GetByUserID(userIdUUID)
}

func (srv *cartServices) Update(data *models.Cart) error{
	dataFromDb, err := srv.cartRepo.GetByID(data.ID)
	if err != nil {
		return err
	}
	if dataFromDb.LinkFile.Valid && dataFromDb.LinkFile != data.LinkFile {
		err := srv.mediaRepo.DeleteFile(context.Background(), dataFromDb.LinkFile.ValueOrZero())
		if err != nil {
			return err
		}
	}

	err = srv.cartRepo.Update(data)
	if err != nil {
		return err
	}
	return nil
}

func (srv *cartServices) DeleteByID(id uint) error{
	return srv.cartRepo.DeleteByCartID(id)
}

func (srv *cartServices) DeleteByUserID(userid string) error{
	userIdUUID, err:= uuid.FromString(userid)
	if err != nil {
		return err
	}
	return srv.cartRepo.DeleteByUserID(userIdUUID)
}
