package usecase

import (
	"log/slog"

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
	logger                *slog.Logger
}

func New(productRepository productRep.RepositoryI, receptionRepository receptionRep.RepositoryI,
	pickupPointRepository pickupPointRep.RepositoryI, logger *slog.Logger) UseCaseI {
	return &useCase{
		productRepository:     productRepository,
		receptionRepository:   receptionRepository,
		pickupPointRepository: pickupPointRepository,
		logger:                logger,
	}
}

func (uc *useCase) CreateProduct(pickupPointID, typeProduct string) (*models.Product, error) {
	err := uc.pickupPointRepository.GetPickupPointByID(pickupPointID)
	if err != nil {
		uc.logger.Error("error with the PP search", slog.String("error", err.Error()))
		return nil, err
	}

	product := &models.Product{}

	product.ReceptionID, _, err = uc.receptionRepository.GetOpenReceptionByPPID(pickupPointID)
	if err != nil {
		uc.logger.Error("error in the database request", slog.String("error", err.Error()))
		return nil, err
	}

	if product.ReceptionID == "" {
		uc.logger.Error("error: there are no open receptions")
		return nil, models.ErrNotOpenReception
	}

	if typeProduct == "электроника" || typeProduct == "одежда" || typeProduct == "обувь" {
		product.TypeProduct = typeProduct
	} else {
		uc.logger.Error("invalid type of product:", slog.String("Entered product type", typeProduct))
		return nil, models.ErrBadData
	}

	err = uc.productRepository.CreateProduct(product)
	if err != nil {
		uc.logger.Error("error in the database request", slog.String("error", err.Error()))
		return nil, err
	}

	return product, nil
}

func (uc *useCase) DeleteLastProduct(pickupPointID string) error {
	err := uc.pickupPointRepository.GetPickupPointByID(pickupPointID)
	if err != nil {
		uc.logger.Error("error in the database request", slog.String("error", err.Error()))
		return err
	}

	receptionID, _, err := uc.receptionRepository.GetOpenReceptionByPPID(pickupPointID)
	if err != nil {
		uc.logger.Error("error with the reception search")
		return err
	}

	if receptionID == "" {
		uc.logger.Error("error: there are no open receptions")
		return models.ErrNotOpenReception
	}

	productID, err := uc.productRepository.GetLastProdcutByDate(receptionID)

	if err != nil {
		uc.logger.Error("error with the product search")
		return err
	}

	err = uc.productRepository.DeleteLastProduct(productID)

	if err != nil {
		uc.logger.Error("error with the product remove")
		return err
	}

	return nil
}
