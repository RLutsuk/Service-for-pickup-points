package delivery

import (
	"errors"
	"log/slog"
	"net/http"

	productUC "github.com/RLutsuk/Service-for-pickup-points/app/internal/product/usecase"
	authMW "github.com/RLutsuk/Service-for-pickup-points/app/internal/user/delivery"
	"github.com/RLutsuk/Service-for-pickup-points/app/models"
	"github.com/labstack/echo/v4"
)

type Delivery struct {
	productUC productUC.UseCaseI
	logger    *slog.Logger
}

func NewDelivery(e *echo.Echo, productUC productUC.UseCaseI, logger *slog.Logger) {
	handler := &Delivery{
		productUC: productUC,
		logger:    logger,
	}
	e.POST("/products", handler.createProduct, authMW.AuthWithRole("employee"))
	e.POST("/pvz/:pvzId/delete_last_product", handler.deleteLastProduct, authMW.AuthWithRole("employee"))
}

// @Summary     Добавление продукта
// @Security    ApiKeyAuth
// @Description Добавление продукта в открытую приемку в выбранном ПВЗ (только для сотрудников ПВЗ)
// @Tags        product
// @Accept      json
// @Produce     json
// @Param       input  body      models.InProduct   true  "Данные продукта"
// @Success     201    {object}  models.OutProduct  	  "Добавленный продукт"
// @Failure     400    {object}  models.ErrorResponse 	  "Неверный запрос"
// @Failure     500    {object}  models.ErrorResponse 	  "Внутрення ошибка сервера"
// @Router      /products [post]
func (delivery *Delivery) createProduct(c echo.Context) error {

	delivery.logger.Info("Request to create a product by user", c.Get("userID"), c.Get("userRole"))

	var body map[string]string
	if err := c.Bind(&body); err != nil {
		delivery.logger.Error("json parsing error, createProduct()", slog.String("error", err.Error()))
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrPickupPointDontExist.Error())
	}

	pickupPointID := body["pvzId"]
	typeProduct := body["type"]

	product, err := delivery.productUC.CreateProduct(pickupPointID, typeProduct)

	if err != nil {
		delivery.logger.Error("error in creating product, createProduct()", slog.String("error", err.Error()))
		return handleProductError(err)
	}

	delivery.logger.Info("successful creation the product", slog.String("product_id", product.ID))
	return c.JSON(http.StatusCreated, product)
}

// @Summary     Удаление продукта
// @Security    ApiKeyAuth
// @Description Удаление продукта из открытой приемки в выбранном ПВЗ (только для сотрудников ПВЗ)
// @Tags        product
// @Accept      json
// @Produce     json
// @Param       pvzId  path      string  				true  "ID ПВЗ"
// @Success     200    {object}  string  					  "Продукт успешно удален"
// @Failure     400    {object}  models.ErrorResponse 	      "Неверный запрос"
// @Failure     500    {object}  models.ErrorResponse 		  "Внутрення ошибка сервера"
// @Router      /pvz/{pvzId}/delete_last_product [post]
func (delivery *Delivery) deleteLastProduct(c echo.Context) error {

	delivery.logger.Info("Request to delete a last product by user", c.Get("userID"), c.Get("userRole"))

	pickupPointID := c.Param("pvzId")

	err := delivery.productUC.DeleteLastProduct(pickupPointID)

	if err != nil {
		delivery.logger.Error("deleteLastProduct()", slog.String("error", err.Error()))
		return handleProductError(err)
	}

	delivery.logger.Info("successful product removal", slog.String("pickup point ID", pickupPointID))
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
