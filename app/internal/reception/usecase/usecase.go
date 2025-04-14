package usecase

import (
	"errors"
	"log/slog"
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
	logger                *slog.Logger
}

func New(receptionRepository receptionRep.RepositoryI, pickupPointRepository pickupPointRep.RepositoryI, logger *slog.Logger) UseCaseI {
	return &useCase{
		receptionRepository:   receptionRepository,
		pickupPointRepository: pickupPointRepository,
		logger:                logger,
	}
}

func (uc *useCase) CreateReception(reception *models.Reception) error {

	err := uc.pickupPointRepository.GetPickupPointByID(reception.PickupPointID)
	if err != nil {
		uc.logger.Error("Error with the PP search", slog.String("error:", err.Error()))
		if errors.Is(err, models.ErrPickupPointDontExist) {
			return models.ErrPickupPointDontExist
		}
		return models.ErrInternalServer
	}

	id, _, err := uc.receptionRepository.GetOpenReceptionByPPID(reception.PickupPointID)
	if err != nil {
		uc.logger.Error("Error with the open reception search", slog.String("error:", err.Error()))
		return err
	}
	if id != "" {
		return models.ErrNotClosedReception
	}

	reception.Status = "in_progress"
	err = uc.receptionRepository.CreateReception(reception)
	if err != nil {
		uc.logger.Error("Error with the creation product", slog.String("error:", err.Error()))
		return err
	}

	return nil
}

func (uc *useCase) CloseReception(pickupPointID string) (*models.Reception, error) {
	err := uc.pickupPointRepository.GetPickupPointByID(pickupPointID)
	if err != nil {
		uc.logger.Error("Error with the PP search", slog.String("error:", err.Error()))
		if errors.Is(err, models.ErrPickupPointDontExist) {
			return nil, err
		}
		return nil, models.ErrBadData
	}

	newReception := &models.Reception{}

	id, dateTime, err := uc.receptionRepository.GetOpenReceptionByPPID(pickupPointID)
	if err != nil {
		uc.logger.Error("Error with the open reception search", slog.String("error:", err.Error()))
		return nil, err
	}
	if id == "" {
		return nil, models.ErrNotOpenReception
	}

	err = uc.receptionRepository.CloseReception(id)
	if err != nil {
		uc.logger.Error("Error with the close reception", slog.String("error:", err.Error()))
		return nil, err
	}

	newReception.DateTime, err = time.Parse(time.RFC3339, dateTime)
	if err != nil {
		uc.logger.Error("Error with the parsing time", slog.String("error:", err.Error()))
		return nil, err
	}

	newReception.ID = id
	newReception.PickupPointID = pickupPointID
	newReception.Status = "close"

	return newReception, nil
}
