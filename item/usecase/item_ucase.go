package usecase

import (
	"context"
	"time"

	"github.com/shop/item"
	"github.com/shop/models"
	//"golang.org/x/sync/errgroup"
)

type itemUsecase struct {
	itemRepo    item.Repository
	contextTimeout time.Duration
}

func NewItemUsecase(a item.Repository, timeout time.Duration) item.Usecase {
	return &itemUsecase{
		itemRepo:    a,
		contextTimeout: timeout,
	}
}


func (a *itemUsecase) Fetch(c context.Context, num int64) ([]*models.Item, error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	listItem, err := a.itemRepo.Fetch(ctx, num)

	if err != nil {
		return nil, err
	}

	return listItem, nil
}

func (a *itemUsecase) Calculate(c context.Context, num int64) (*models.Calculate, error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	calculateItem, err := a.itemRepo.Calculate(ctx, num)

	if err != nil {
		return nil, err
	}

	return calculateItem, nil
}


func (a *itemUsecase) Store(c context.Context, m *models.Item) error {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	err := a.itemRepo.Store(ctx, m)
	if err != nil {
		return err
	}
	return nil
}

