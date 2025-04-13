package delivery

import (
	"errors"
	"net/http"

	userUC "github.com/RLutsuk/Service-for-pickup-points/app/internal/user/usecase"
	"github.com/RLutsuk/Service-for-pickup-points/app/models"
	"github.com/labstack/echo"
)

type Delivery struct {
	userUC userUC.UseCaseI
}

func NewDelivery(e *echo.Echo, userUC userUC.UseCaseI) {
	handler := &Delivery{
		userUC: userUC,
	}
	e.POST("/register", handler.createUser)
	e.POST("/login", handler.authUser)
	e.POST("/dummyLogin", handler.testToken)
}

func (delivery *Delivery) createUser(c echo.Context) error {

	var user models.User
	if err := c.Bind(&user); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData)
	}

	err := delivery.userUC.CreateUser(&user)

	if err != nil {
		c.Logger().Error(err)
		return handleUserError(err)
	}

	return c.JSON(http.StatusOK, user)
}

func (delivery *Delivery) authUser(c echo.Context) error {

	var user models.User
	if err := c.Bind(&user); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData.Error())
	}

	token, err := delivery.userUC.AuthUser(&user)

	if err != nil {
		c.Logger().Error(err)
		return handleUserError(err)
	}

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
		c.Logger().Error(err)
		return handleUserError(err)
	}

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
