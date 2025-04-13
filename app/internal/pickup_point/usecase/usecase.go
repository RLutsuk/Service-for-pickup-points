package usecase

import (
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
}

func New(pickupPointRepository pickupPointRep.RepositoryI) UseCaseI {
	return &useCase{
		pickupPointRepository: pickupPointRepository,
	}
}

func (uc *useCase) CreatePickupPoint(pickupPoint *models.PickupPoint) error {

	if pickupPoint.City != "Москва" && pickupPoint.City != "Санкт-Петербург" &&
		pickupPoint.City != "Казань" {
		return models.ErrBadData
	}

	err := uc.pickupPointRepository.CreatePickupPoint(pickupPoint)

	if err != nil {
		return err
	}

	return nil
}
func (uc *useCase) GetAllPickupPoint(startDate, endDate, page, limit string) ([]*models.PickupPoint, error) {

	limitInt, err := strconv.Atoi(limit)
	if limit != "" && err != nil {
		return nil, models.ErrBadData
	}

	if limitInt < 1 {
		limitInt = 5
	}

	offset, err := strconv.Atoi(page)
	if page != "" && err != nil {
		return nil, models.ErrBadData
	}

	offset = (offset - 1) * limitInt

	if offset < 0 {
		offset = 0
	}

	if startDate != "" {
		_, err = time.Parse(time.RFC3339, startDate)
		if err != nil {
			return nil, models.ErrBadData
		}
	}

	if endDate != "" {
		_, err = time.Parse(time.RFC3339, endDate)
		if err != nil {
			return nil, models.ErrBadData
		}
	}

	pickupPoints, err := uc.pickupPointRepository.GetAllPickupPoint(startDate, endDate, offset, limitInt)

	if err != nil {
		return nil, err
	}

	return pickupPoints, nil
}
