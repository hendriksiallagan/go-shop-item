package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	//"strconv"
	"strings"
	"testing"
	//"time"

	"github.com/labstack/echo"
	itemHttp "github.com/go-shop-item/item/delivery/http"
	"github.com/go-shop-item/item/mocks"
	"github.com/go-shop-item/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/bxcodec/faker"
)

func TestFetch(t *testing.T) {
	var mockItem models.Item
	err := faker.FakeData(&mockItem)
	assert.NoError(t, err)
	mockICase := new(mocks.Usecase)
	mockListItem := make([]*models.Item, 0)
	mockListItem = append(mockListItem, &mockItem)
	num := 1
	cursor := "2"
	mockICase.On("Fetch", mock.Anything, cursor, int64(num)).Return(mockListItem, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/item", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := itemHttp.HttpItemHandler{
		IUsecase: mockICase,
	}
	handler.FetchItem(c)

	responseCursor := rec.Header().Get("X-Cursor")
	assert.Equal(t, "10", responseCursor)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockICase.AssertExpectations(t)
}

func TestFetchError(t *testing.T) {
	mockUCase := new(mocks.Usecase)
	num := 1
	cursor := "2"
	mockUCase.On("Fetch", mock.Anything, cursor, int64(num)).Return(nil, "", models.ErrInternalServerError)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/item", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := itemHttp.HttpItemHandler{
		IUsecase: mockUCase,
	}
	handler.FetchItem(c)

	responseCursor := rec.Header().Get("X-Cursor")
	assert.Equal(t, "", responseCursor)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestStore(t *testing.T) {
	mockItem := models.Item{
		Name:     "Big Mac",
		Price:   1000,
		TaxCode: 1,
	}

	tempMockItem := mockItem
	tempMockItem.ID = 0
	mockUCase := new(mocks.Usecase)

	j, err := json.Marshal(tempMockItem)
	assert.NoError(t, err)

	mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*models.Item")).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/items", strings.NewReader(string(j)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/items")

	handler := itemHttp.HttpItemHandler{
		IUsecase: mockUCase,
	}
	handler.Store(c)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestCalculate(t *testing.T) {
	var mockCalculate models.Calculate
	err := faker.FakeData(&mockCalculate)
	assert.NoError(t, err)

	mockUCase := new(mocks.Usecase)
	mockListCalculate := make([]*models.Calculate, 0)
	mockListCalculate = append(mockListCalculate, &mockCalculate)
	num := 1
	mockUCase.On("Calculate", mock.Anything, int64(num)).Return(mockListCalculate, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/item", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := itemHttp.HttpItemHandler{
		IUsecase: mockUCase,
	}
	handler.CalculateItem(c)

	responseCursor := rec.Header().Get("X-Cursor")
	assert.Equal(t, "10", responseCursor)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

