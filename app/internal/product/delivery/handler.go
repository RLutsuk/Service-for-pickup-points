package delivery

import (
	"errors"
	"net/http"

	productUC "github.com/RLutsuk/Service-for-pickup-points/app/internal/product/usecase"
	authMW "github.com/RLutsuk/Service-for-pickup-points/app/internal/user/delivery"
	"github.com/RLutsuk/Service-for-pickup-points/app/models"
	"github.com/labstack/echo"
)

type Delivery struct {
	productUC productUC.UseCaseI
}

func NewDelivery(e *echo.Echo, productUC productUC.UseCaseI) {
	handler := &Delivery{
		productUC: productUC,
	}
	e.POST("/products", handler.createProduct, authMW.AuthWithRole("employee"))
	e.POST("/pvz/:pvzId/delete_last_product", handler.deleteLastProduct, authMW.AuthWithRole("employee"))
}

func (delivery *Delivery) createProduct(c echo.Context) error {

	var body map[string]string
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrPickupPointDontExist.Error())
	}

	pickupPointID := body["pvzId"]
	typeProduct := body["type"]

	product, err := delivery.productUC.CreateProduct(pickupPointID, typeProduct)

	if err != nil {
		c.Logger().Error(err)
		return handleProductError(err)
	}

	return c.JSON(http.StatusOK, product)
}

func (delivery *Delivery) deleteLastProduct(c echo.Context) error {

	pickupPointID := c.Param("pvzId")

	err := delivery.productUC.DeleteLastProduct(pickupPointID)

	if err != nil {
		c.Logger().Error(err)
		return handleProductError(err)
	}

	return c.JSON(http.StatusOK, "The product was successfully deleted")
}

func handleProductError(err error) *echo.HTTPError {
	switch {
	case errors.Is(err, models.ErrBadData):
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData.Error())
	case errors.Is(err, models.ErrNotOpenReception):
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrNotOpenReception.Error())
	case errors.Is(err, models.ErrPickupPointDontExist):
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrPickupPointDontExist.Error())
	case errors.Is(err, models.ErrEmptyReception):
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrEmptyReception.Error())
	}
	return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServer.Error())
}
