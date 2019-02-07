package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-shop-item/models"
	"github.com/sirupsen/logrus"

	"github.com/labstack/echo"
	"github.com/go-shop-item/item"

	"gopkg.in/go-playground/validator.v9"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// HttpItemHandler  represent the httphandler for article
type HttpItemHandler struct {
	IUsecase item.Usecase
}

func NewArticleHttpHandler(e *echo.Echo, us item.Usecase) {
	handler := &HttpItemHandler{
		IUsecase: us,
	}
	e.GET("/items", handler.FetchItem)
	e.GET("/calculate", handler.CalculateItem)
	e.POST("/items", handler.Store)
}

func (a *HttpItemHandler) FetchItem(c echo.Context) error {
	numS := c.QueryParam("num")
	num, _ := strconv.Atoi(numS)
	cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	listAr,  err := a.IUsecase.Fetch(ctx, int64(num))

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	c.Response().Header().Set(`X-Cursor`, cursor)

	return c.JSON(http.StatusOK, listAr)
}

func (a *HttpItemHandler) CalculateItem(c echo.Context) error {
	numS := c.QueryParam("num")
	num, _ := strconv.Atoi(numS)
	cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	listAr,  err := a.IUsecase.Calculate(ctx, int64(num))

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	c.Response().Header().Set(`X-Cursor`, cursor)

	return c.JSON(http.StatusOK, listAr)
}

func isRequestValid(m *models.Item) (bool, error) {

	validate := validator.New()

	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *HttpItemHandler) Store(c echo.Context) error {
	var item models.Item
	err := c.Bind(&item)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&item); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = a.IUsecase.Store(ctx, &item)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, item)
}


func getStatusCode(err error) int {

	if err == nil {
		return http.StatusOK
	}
	logrus.Error(err)
	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
