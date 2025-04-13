package delivery

import (
	"errors"
	"net/http"

	pickupPointUC "github.com/RLutsuk/Service-for-pickup-points/app/internal/pickup_point/usecase"
	authMW "github.com/RLutsuk/Service-for-pickup-points/app/internal/user/delivery"
	"github.com/RLutsuk/Service-for-pickup-points/app/models"
	"github.com/labstack/echo"
)

type Delivery struct {
	pickupPointUC pickupPointUC.UseCaseI
}

func NewDelivery(e *echo.Echo, pickupPointUC pickupPointUC.UseCaseI) {
	handler := &Delivery{
		pickupPointUC: pickupPointUC,
	}
	e.POST("/pvz", handler.createPickupPoint, authMW.AuthWithRole("moderator"))
	e.GET("/pvz", handler.getAllPickupPoints, authMW.AuthWithRole("employee", "moderator"))
}

func (delivery *Delivery) createPickupPoint(c echo.Context) error {
	var pickupPoint models.PickupPoint

	if err := c.Bind(&pickupPoint); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData.Error())
	}

	err := delivery.pickupPointUC.CreatePickupPoint(&pickupPoint)

	if err != nil {
		c.Logger().Error(err)
		return handlePickupPointError(err)
	}

	return c.JSON(http.StatusOK, pickupPoint)
}

func (delivery *Delivery) getAllPickupPoints(c echo.Context) error {
	startDate := c.QueryParam("startDate")
	endDate := c.QueryParam("endDate")
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")

	pickupPoints, err := delivery.pickupPointUC.GetAllPickupPoint(startDate, endDate, page, limit)

	if err != nil {
		c.Logger().Error(err)
		return handlePickupPointError(err)
	}

	return c.JSON(http.StatusOK, pickupPoints)
}

func handlePickupPointError(err error) *echo.HTTPError {
	switch {
	case errors.Is(err, models.ErrBadData):
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData.Error())
	}
	return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServer.Error())
}
