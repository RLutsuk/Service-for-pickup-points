package usecase

import (
	"log/slog"
	"strconv"
	"time"

	pickupPointRep "github.com/RLutsuk/Service-for-pickup-points/app/internal/pickup_point/repository"
	"github.com/RLutsuk/Service-for-pickup-points/app/models"
)

type UseCaseI interface {
	CreatePickupPoint(pickupPoint *models.PickupPoint) error
	GetAllPickupPoint(startDate, endDate, page, limit string) ([]*models.PickupPoint, error)
}

type useCase struct {
	pickupPointRepository pickupPointRep.RepositoryI
	logger                *slog.Logger
}

func New(pickupPointRepository pickupPointRep.RepositoryI, logger *slog.Logger) UseCaseI {
	return &useCase{
		pickupPointRepository: pickupPointRepository,
		logger:                logger,
	}
}

func (uc *useCase) CreatePickupPoint(pickupPoint *models.PickupPoint) error {

	if pickupPoint.City != "Москва" && pickupPoint.City != "Санкт-Петербург" &&
		pickupPoint.City != "Казань" {
			uc.logger.Error("Invalid city")
			return models.ErrBadData
	}

	uc.logger.Debug("Database request")
	err := uc.pickupPointRepository.CreatePickupPoint(pickupPoint)

	if err != nil {
		uc.logger.Error("Error in the database request")
		return err
	}

	return nil
}

func (uc *useCase) GetAllPickupPoint(startDate, endDate, page, limit string) ([]*models.PickupPoint, error) {

	limitInt, err := strconv.Atoi(limit)
	if limit != "" && err != nil {
		uc.logger.Error("Invalid limit")
		return nil, models.ErrBadData
	}

	if limitInt < 1 || limitInt > 25 {
		limitInt = 5
	}

	offset, err := strconv.Atoi(page)
	if page != "" && err != nil {
		uc.logger.Error("Invalid offset")
		return nil, models.ErrBadData
	}

	offset = (offset - 1) * limitInt

	if offset < 0 {
		offset = 0
	}

	if startDate != "" {
		_, err = time.Parse(time.RFC3339, startDate)
		if err != nil {
			uc.logger.Error("Invalid startDate")
			return nil, models.ErrBadData
		}
	}

	if endDate != "" {
		_, err = time.Parse(time.RFC3339, endDate)
		if err != nil {
			uc.logger.Error("Invalid endDate")
			return nil, models.ErrBadData
		}
	}

	uc.logger.Debug("Database request")
	pickupPoints, err := uc.pickupPointRepository.GetAllPickupPoint(startDate, endDate, offset, limitInt)

	if err != nil {
		uc.logger.Error("Error in the database request")
		return nil, err
	}

	return pickupPoints, nil
}
