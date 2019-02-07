package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-shop-item/item/mocks"
	ucase "github.com/go-shop-item/item/usecase"
	"github.com/go-shop-item/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetch(t *testing.T) {
	mockItemRepo := new(mocks.Repository)
	mockItem := &models.Item{
		Name:   "Big Mac",
		Price: 1000,
		TaxCode: 1,
	}

	mockListItem := make([]*models.Item, 0)
	mockListItem = append(mockListItem, mockItem)

	t.Run("success", func(t *testing.T) {
		mockItemRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(mockListItem,  nil).Once()

		u := ucase.NewItemUsecase(mockItemRepo, time.Second*2)
		num := int64(1)
		list, err := u.Fetch(context.TODO(), num)


		assert.NoError(t, err)
		assert.Len(t, list, len(mockListItem))

		mockItemRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockItemRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(nil, "", errors.New("Unexpexted Error")).Once()

		u := ucase.NewItemUsecase(mockItemRepo, time.Second*2)
		num := int64(1)
		list, err := u.Fetch(context.TODO(), num)

		assert.Error(t, err)
		assert.Len(t, list, 0)
		mockItemRepo.AssertExpectations(t)
	})

}

func TestCalculate(t *testing.T) {
	mockItemRepo := new(mocks.Repository)
	mockItem := models.Calculate{
		PriceSubtotal: 2150,
		TaxSubtotal: 135,
		GrandTotal: 2285,
	}

	t.Run("success", func(t *testing.T) {
		mockItemRepo.On("Calculate", mock.Anything, mock.AnythingOfType("int64")).Return(&mockItem, nil).Once()

		u := ucase.NewItemUsecase(mockItemRepo, time.Second*2)

		a, err := u.Calculate(context.TODO(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockItemRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockItemRepo.On("Calculate", mock.Anything, mock.AnythingOfType("int64")).Return(nil, errors.New("Unexpected")).Once()
		u := ucase.NewItemUsecase(mockItemRepo, time.Second*2)

		a, err := u.Calculate(context.TODO(),1)

		assert.Error(t, err)
		assert.Nil(t, a)

		mockItemRepo.AssertExpectations(t)
	})

}

func TestStore(t *testing.T) {
	mockItemRepo := new(mocks.Repository)
	mockItem := models.Item{
		Name:   "Big Mac",
		Price: 1000,
		TaxCode: 1,
	}

	t.Run("success", func(t *testing.T) {
		tempMockItem := mockItem
		tempMockItem.ID = 0
		mockItemRepo.On("GetByTitle", mock.Anything, mock.AnythingOfType("string")).Return(nil, models.ErrNotFound).Once()
		mockItemRepo.On("Store", mock.Anything, mock.AnythingOfType("*models.Article")).Return(nil).Once()

		u := ucase.NewItemUsecase(mockItemRepo, time.Second*2)

		err := u.Store(context.TODO(), &tempMockItem)

		assert.NoError(t, err)
		assert.Equal(t, mockItem.Name, tempMockItem.Name)
		mockItemRepo.AssertExpectations(t)
	})

}
