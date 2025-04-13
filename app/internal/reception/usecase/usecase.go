package usecase

import (
	"errors"
	"time"

	pickupPointRep "github.com/RLutsuk/Service-for-pickup-points/app/internal/pickup_point/repository"
	receptionRep "github.com/RLutsuk/Service-for-pickup-points/app/internal/reception/repository"
	"github.com/RLutsuk/Service-for-pickup-points/app/models"
)

type UseCaseI interface {
	CreateReception(reception *models.Reception) error
	CloseReception(pickupPointID string) (*models.Reception, error)
}

type useCase struct {
	receptionRepository   receptionRep.RepositoryI
	pickupPointRepository pickupPointRep.RepositoryI
}

func New(receptionRepository receptionRep.RepositoryI, pickupPointRepository pickupPointRep.RepositoryI) UseCaseI {
	return &useCase{
		receptionRepository:   receptionRepository,
		pickupPointRepository: pickupPointRepository,
	}
}

func (uc *useCase) CreateReception(reception *models.Reception) error {

	err := uc.pickupPointRepository.GetPickupPointByID(reception.PickupPointID)
	if err != nil {
		if errors.Is(err, models.ErrPickupPointDontExist) {
			return models.ErrPickupPointDontExist
		}
		return models.ErrInternalServer
	}

	id, _, err := uc.receptionRepository.GetOpenReceptionByPPID(reception.PickupPointID)
	if err != nil {
		return err
	}
	if id != "" {
		return models.ErrNotClosedReception
	}

	reception.Status = "in_progress"
	err = uc.receptionRepository.CreateReception(reception)
	if err != nil {
		return err
	}

	return nil
}

func (uc *useCase) CloseReception(pickupPointID string) (*models.Reception, error) {
	err := uc.pickupPointRepository.GetPickupPointByID(pickupPointID)
	if err != nil {
		if errors.Is(err, models.ErrPickupPointDontExist) {
			return nil, err
		}
		return nil, models.ErrBadData
	}

	newReception := &models.Reception{}

	id, dateTime, err := uc.receptionRepository.GetOpenReceptionByPPID(pickupPointID)
	if err != nil {
		return nil, err
	}
	if id == "" {
		return nil, models.ErrNotOpenReception
	}

	err = uc.receptionRepository.CloseReception(id)
	if err != nil {
		return nil, err
	}

	newReception.DateTime, err = time.Parse(time.RFC3339, dateTime)
	if err != nil {
		return nil, err
	}

	newReception.ID = id
	newReception.PickupPointID = pickupPointID
	newReception.Status = "close"

	return newReception, nil
}
