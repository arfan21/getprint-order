package services

import (
	"context"

	"github.com/arfan21/getprint-order/app/models"
	_mediaRepo "github.com/arfan21/getprint-order/app/repository/media"
	repo "github.com/arfan21/getprint-order/app/repository/mysql"
	"github.com/arfan21/getprint-order/configs"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/sync/errgroup"
)

type CartServices interface {
	Create(data []models.Cart) error
	GetByUserID(userid string) (*[]models.Cart, error)
	UpdateBatch(data []models.Cart) error
	DeleteByID(id uint) error
	DeleteByUserID(userid string) error
}

type cartServices struct {
	cartRepo  repo.CartRepository
	mediaRepo _mediaRepo.MediaRepository
	db        configs.Client
}

func NewCartServices(db configs.Client, cartRepo repo.CartRepository) CartServices {
	mediaRepo := _mediaRepo.NewMediaRepository()
	return &cartServices{cartRepo, mediaRepo, db}
}

func (srv *cartServices) Create(data []models.Cart) error {
	userId := data[0].UserID
	dataCarts, err := srv.cartRepo.GetByUserID(userId)
	if err != nil {
		return err
	}

	for _, cart := range *dataCarts {
		for i, newCart := range data {
			if cart.OrderType == newCart.OrderType && cart.PartnerID == newCart.PartnerID {
				cart.Qty += newCart.Qty
				err := srv.cartRepo.UpdateByID(&cart)
				if err != nil {
					return err
				}

				// delete cart
				data[i] = data[len(data)-1]
				data[len(data)-1] = models.Cart{}
				data = data[:len(data)-1]
			}
		}
	}

	if len(data) == 0 {
		return nil
	}

	return srv.cartRepo.Create(data)
}

func (srv *cartServices) GetByUserID(userid string) (*[]models.Cart, error) {
	userIdUUID, err := uuid.FromString(userid)
	if err != nil {
		return nil, err
	}
	return srv.cartRepo.GetByUserID(userIdUUID)
}

func (srv *cartServices) UpdateBatch(data []models.Cart) error {
	errg := errgroup.Group{}
	tx := srv.db.Conn().Begin()

	for _, cartItem := range data {
		newCartItem := cartItem
		errg.Go(func() error {
			err := srv.cartRepo.UpdateByIDwithTx(tx, newCartItem)
			if err != nil {
				return err
			}
			return nil
		})
	}

	if err := errg.Wait(); err != nil {
		return err
	}

	return tx.Commit().Error
}

func (srv *cartServices) DeleteByID(id uint) error {
	dataFromDb, err := srv.cartRepo.GetByID(id)
	if err != nil {
		return err
	}

	if dataFromDb.LinkFile.Valid {
		err := srv.mediaRepo.DeleteFile(context.Background(), dataFromDb.LinkFile.ValueOrZero())
		if err != nil {
			return err
		}
	}

	return srv.cartRepo.DeleteByCartID(id)
}

func (srv *cartServices) DeleteByUserID(userid string) error {
	userIdUUID, err := uuid.FromString(userid)
	if err != nil {
		return err
	}

	dataFromDb, err := srv.cartRepo.GetByUserID(userIdUUID)
	if err != nil {
		return err
	}

	errg, ctx := errgroup.WithContext(context.Background())

	for _, data := range *dataFromDb {
		if data.LinkFile.Valid {
			errg.Go(func() error {
				err := srv.mediaRepo.DeleteFile(ctx, data.LinkFile.ValueOrZero())
				if err != nil {
					return err
				}

				return nil
			})

		}
	}

	if err := errg.Wait(); err != nil {
		return err
	}

	return srv.cartRepo.DeleteByUserID(userIdUUID)
}
