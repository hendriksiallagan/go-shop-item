package item

import (
	"context"
	"github.com/shop/models"
)

// Usecase represent the item's usecases
type Usecase interface {
	Fetch(ctx context.Context, num int64) ([]*models.Item, error)
	Calculate(ctx context.Context, num int64) (*models.Calculate, error)
	Store(context.Context, *models.Item) error
}
