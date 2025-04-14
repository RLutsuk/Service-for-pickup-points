package delivery

import (
	"errors"
	"log/slog"
	"net/http"

	userUC "github.com/RLutsuk/Service-for-pickup-points/app/internal/user/usecase"
	"github.com/RLutsuk/Service-for-pickup-points/app/models"
	"github.com/labstack/echo"
)

type Delivery struct {
	userUC userUC.UseCaseI
	logger *slog.Logger
}

func NewDelivery(e *echo.Echo, userUC userUC.UseCaseI, logger *slog.Logger) {
	handler := &Delivery{
		userUC: userUC,
		logger: logger,
	}
	e.POST("/register", handler.createUser)
	e.POST("/login", handler.authUser)
	e.POST("/dummyLogin", handler.testToken)
}

func (delivery *Delivery) createUser(c echo.Context) error {

	delivery.logger.Info("Request to create user")

	var user models.User
	if err := c.Bind(&user); err != nil {
		delivery.logger.Error("json parsing error, createUser()", slog.String("error", err.Error()))
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData)
	}

	err := delivery.userUC.CreateUser(&user)

	if err != nil {
		delivery.logger.Error("createUser()", slog.String("error", err.Error()))
		return handleUserError(err)
	}

	delivery.logger.Info("successful creation the user", slog.String("user ID", user.ID))
	return c.JSON(http.StatusCreated, user)
}

func (delivery *Delivery) authUser(c echo.Context) error {

	var user models.User
	if err := c.Bind(&user); err != nil {
		delivery.logger.Error("authUser()", slog.String("error", err.Error()))
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData.Error())
	}

	token, err := delivery.userUC.AuthUser(&user)

	if err != nil {
		delivery.logger.Error("authUser()", slog.String("error", err.Error()))
		return handleUserError(err)
	}

	delivery.logger.Info("successful auth of user", slog.String("userID", user.ID))
	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func (delivery *Delivery) testToken(c echo.Context) error {

	var body map[string]string
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrPickupPointDontExist.Error())
	}

	role := body["role"]

	token, err := delivery.userUC.TestUser(role)

	if err != nil {
		delivery.logger.Error("authUser()", slog.String("error", err.Error()))
		return handleUserError(err)
	}

	delivery.logger.Info("token was successfully issued")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func handleUserError(err error) *echo.HTTPError {
	switch {
	case errors.Is(err, models.ErrBadData):
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData.Error())
	case errors.Is(err, models.ErrBadEmail):
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadEmail.Error())
	}
	return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServer.Error())
}
