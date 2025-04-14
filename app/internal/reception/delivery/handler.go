package delivery

import (
	"errors"
	"log/slog"
	"net/http"

	receptionUC "github.com/RLutsuk/Service-for-pickup-points/app/internal/reception/usecase"
	authMW "github.com/RLutsuk/Service-for-pickup-points/app/internal/user/delivery"
	"github.com/RLutsuk/Service-for-pickup-points/app/models"
	"github.com/labstack/echo"
)

type Delivery struct {
	receptionUC receptionUC.UseCaseI
	logger      *slog.Logger
}

func NewDelivery(e *echo.Echo, receptionUC receptionUC.UseCaseI, logger *slog.Logger) {
	handler := &Delivery{
		receptionUC: receptionUC,
		logger:      logger,
	}
	e.POST("/receptions", handler.createReception, authMW.AuthWithRole("employee"))
	e.POST("/pvz/:pvzId/close_last_reception", handler.closeReception, authMW.AuthWithRole("employee"))
}

func (delivery *Delivery) createReception(c echo.Context) error {

	delivery.logger.Info("Request to create a receprion by user", c.Get("userID"), c.Get("userRole"))

	var reception models.Reception

	if err := c.Bind(&reception); err != nil {
		delivery.logger.Error("json parsing error, createReception()", slog.String("error:", err.Error()))
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData)
	}

	err := delivery.receptionUC.CreateReception(&reception)

	if err != nil {
		delivery.logger.Error("createReception()", slog.String("error:", err.Error()))
		return handleReceptionError(err)
	}

	delivery.logger.Info("successful creation the reception", slog.String("reception_id", reception.ID))
	return c.JSON(http.StatusCreated, reception)
}

func (delivery *Delivery) closeReception(c echo.Context) error {

	delivery.logger.Info("Request to close a receprion by user", c.Get("userID"), c.Get("userRole"))

	pickupPointID := c.Param("pvzId")
	reception, err := delivery.receptionUC.CloseReception(pickupPointID)

	if err != nil {
		delivery.logger.Error("closeReception()", slog.String("error:", err.Error()))
		return handleReceptionError(err)
	}

	delivery.logger.Info("successful reception closure", slog.String("reception_id", reception.ID))
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
