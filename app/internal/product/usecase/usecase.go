package usecase

import (
	pickupPointRep "github.com/RLutsuk/Service-for-pickup-points/app/internal/pickup_point/repository"
	productRep "github.com/RLutsuk/Service-for-pickup-points/app/internal/product/repository"
	receptionRep "github.com/RLutsuk/Service-for-pickup-points/app/internal/reception/repository"
	"github.com/RLutsuk/Service-for-pickup-points/app/models"
)

type UseCaseI interface {
	CreateProduct(pickupPointID, typeProduct string) (*models.Product, error)
	DeleteLastProduct(pickupPointID string) error
}

type useCase struct {
	productRepository     productRep.RepositoryI
	receptionRepository   receptionRep.RepositoryI
	pickupPointRepository pickupPointRep.RepositoryI
}

func New(productRepository productRep.RepositoryI, receptionRepository receptionRep.RepositoryI,
	pickupPointRepository pickupPointRep.RepositoryI) UseCaseI {
	return &useCase{
		productRepository:     productRepository,
		receptionRepository:   receptionRepository,
		pickupPointRepository: pickupPointRepository,
	}
}

func (uc *useCase) CreateProduct(pickupPointID, typeProduct string) (*models.Product, error) {
	err := uc.pickupPointRepository.GetPickupPointByID(pickupPointID)
	if err != nil {
		return nil, err
	}

	product := &models.Product{}

	product.ReceptionID, _, err = uc.receptionRepository.GetOpenReceptionByPPID(pickupPointID)
	if err != nil {
		return nil, err
	}

	if product.ReceptionID == "" {
		return nil, models.ErrNotOpenReception
	}

	if typeProduct == "электроника" || typeProduct == "одежда" || typeProduct == "обувь" {
		product.TypeProduct = typeProduct
	} else {
		return nil, models.ErrBadData
	}

	err = uc.productRepository.CreateProduct(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (uc *useCase) DeleteLastProduct(pickupPointID string) error {
	err := uc.pickupPointRepository.GetPickupPointByID(pickupPointID)
	if err != nil {
		return err
	}

	receptionID, _, err := uc.receptionRepository.GetOpenReceptionByPPID(pickupPointID)
	if err != nil {
		return err
	}

	if receptionID == "" {
		return models.ErrNotOpenReception
	}

	productID, err := uc.productRepository.GetLastProdcutByDate(receptionID)

	if err != nil {
		return err
	}

	err = uc.productRepository.DeleteLastProduct(productID)

	if err != nil {
		return err
	}

	return nil
}
