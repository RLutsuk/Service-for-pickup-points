package delivery

import (
	"errors"
	"log/slog"
	"net/http"

	middleware "github.com/RLutsuk/Service-for-pickup-points/app/internal/middleware"
	pickupPointUC "github.com/RLutsuk/Service-for-pickup-points/app/internal/pickup_point/usecase"
	"github.com/RLutsuk/Service-for-pickup-points/app/models"
	"github.com/labstack/echo/v4"
)

type Delivery struct {
	pickupPointUC pickupPointUC.UseCaseI
	logger        *slog.Logger
}

func NewDelivery(e *echo.Echo, pickupPointUC pickupPointUC.UseCaseI, logger *slog.Logger) {
	handler := &Delivery{
		pickupPointUC: pickupPointUC,
		logger:        logger,
	}
	e.POST("/pvz", handler.createPickupPoint, middleware.AuthWithRole("moderator"))
	e.GET("/pvz", handler.getAllPickupPoints, middleware.AuthWithRole("employee", "moderator"))
}

// @Summary     Создание ПВЗ
// @Description Создание ПВЗ в сервисе (только для модераторов)
// @Tags        pvz
// @Security    ApiKeyAuth
// @Accept      json
// @Produce     json
// @Param       input  body      models.CreationInPickupPoint   true  "Данные для ПВЗ"
// @Success     201    {object}  models.CreationOutPickupPoint        "Созданный ПВЗ"
// @Failure     400    {object}  models.ErrorResponse                 "Неверный запрос"
// @Failure     500    {object}  models.ErrorResponse                 "Внутренняя ошибка сервера"
// @Router      /pvz [post]
func (delivery *Delivery) createPickupPoint(c echo.Context) error {

	delivery.logger.Info("Request to create a pickup point by user", c.Get("userID"), c.Get("userRole"))

	var pickupPoint models.PickupPoint

	if err := c.Bind(&pickupPoint); err != nil {
		delivery.logger.Error("json parsing error, createPickupPoint()", slog.String("error", err.Error()))
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData.Error())
	}

	err := delivery.pickupPointUC.CreatePickupPoint(&pickupPoint)

	if err != nil {
		delivery.logger.Error("error in creating pp, createPickupPoint()", slog.String("error", err.Error()))
		return handlePickupPointError(err)
	}

	delivery.logger.Info("successful creation the pickup point", slog.String("pickup_point", pickupPoint.ID))
	return c.JSON(http.StatusCreated, pickupPoint)
}

// @Summary      Получение списка всех ПВЗ
// @Security 	 ApiKeyAuth
// @Description  Получение списка всех ПВЗ с их приемками и товарами (только для модераторов и сотрудника ПВЗ)
// @Tags         pvz
// @Accept       json
// @Produce      json
// @Param        startDate query     string  				false  "startDate"
// @Param        endDate   query     string 			 	false  "endDate"
// @Param        limit     query     int     				false  "limit"  default(5)
// @Param        page      query     int     				false  "page" default(1)
// @Success      200       {array}   models.PickupPoint 		   "Список всех ПВЗ"
// @Failure      400       {object}  models.ErrorResponse 		   "Неверный запрос"
// @Failure      500       {object}  models.ErrorResponse  	       "Внутрення ошибка сервера"
// @Router       /pvz [get]
func (delivery *Delivery) getAllPickupPoints(c echo.Context) error {

	delivery.logger.Info("Request for all PP by user", c.Get("userID"), c.Get("userRole"))

	startDate := c.QueryParam("startDate")
	endDate := c.QueryParam("endDate")
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")

	pickupPoints, err := delivery.pickupPointUC.GetAllPickupPoint(startDate, endDate, page, limit)

	if err != nil {
		delivery.logger.Error("error in request PP, getAllPickupPoints()", slog.String("error", err.Error()))
		return handlePickupPointError(err)
	}

	delivery.logger.Info("All PPs have been successfully found")
	return c.JSON(http.StatusOK, pickupPoints)
}

func handlePickupPointError(err error) *echo.HTTPError {
	switch {
	case errors.Is(err, models.ErrBadData):
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData.Error())
	}
	return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServer.Error())
}
