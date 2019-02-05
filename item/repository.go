package item

import (
	"context"
	"github.com/shop/models"
)

// Repository represent the item's repository contract
type Repository interface {
	Fetch(ctx context.Context, num int64) (res []*models.Item, err error)
	Calculate(ctx context.Context, num int64) (*models.Calculate, error)
	Store(ctx context.Context, a *models.Item) error
}
