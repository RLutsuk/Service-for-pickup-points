package delivery

import (
	"errors"
	"net/http"

	receptionUC "github.com/RLutsuk/Service-for-pickup-points/app/internal/reception/usecase"
	authMW "github.com/RLutsuk/Service-for-pickup-points/app/internal/user/delivery"
	"github.com/RLutsuk/Service-for-pickup-points/app/models"
	"github.com/labstack/echo"
)

type Delivery struct {
	receptionUC receptionUC.UseCaseI
}

func NewDelivery(e *echo.Echo, receptionUC receptionUC.UseCaseI) {
	handler := &Delivery{
		receptionUC: receptionUC,
	}
	e.POST("/receptions", handler.createReception, authMW.AuthWithRole("employee"))
	e.POST("/pvz/:pvzId/close_last_reception", handler.closeReception, authMW.AuthWithRole("employee"))
}

func (delivery *Delivery) createReception(c echo.Context) error {
	var reception models.Reception

	if err := c.Bind(&reception); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData)
	}

	err := delivery.receptionUC.CreateReception(&reception)

	if err != nil {
		c.Logger().Error(err)
		return handleReceptionError(err)
	}

	return c.JSON(http.StatusOK, reception)
}

func (delivery *Delivery) closeReception(c echo.Context) error {

	pickupPointID := c.Param("pvzId")
	reception, err := delivery.receptionUC.CloseReception(pickupPointID)

	if err != nil {
		c.Logger().Error(err)
		return handleReceptionError(err)
	}

	return c.JSON(http.StatusOK, reception)
}

func handleReceptionError(err error) *echo.HTTPError {
	switch {
	case errors.Is(err, models.ErrBadData):
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData.Error())
	case errors.Is(err, models.ErrPickupPointDontExist):
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrPickupPointDontExist.Error())
	case errors.Is(err, models.ErrNotOpenReception):
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrNotOpenReception.Error())
	case errors.Is(err, models.ErrNotClosedReception):
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrNotClosedReception.Error())
	}
	return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServer.Error())
}
