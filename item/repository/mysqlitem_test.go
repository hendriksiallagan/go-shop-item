package repository_test

import (
	"context"
	"testing"

	itemRepo "github.com/shop/item/repository"
	"github.com/shop/models"
	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockItems := []models.Item{
		models.Item{
			ID: 1, Name: "Big Mac", Price: 1000, TaxCode: 1, Type: "Food & Beverage", Refundable: "Refundable", Tax: 100, Amount: 1100,
		},
		models.Item{
			ID: 2, Name: "Lucky Stretch", Price: 1000, TaxCode:2, Type: "Tobacco", Refundable: "Not Refundable", Tax: 30, Amount: 1030,
		},
	}

	rows := sqlmock.NewRows([]string{"id", "name", "price", "tax_code", "type", "refundable", "tax", "amount"}).
		AddRow(mockItems[0].ID, mockItems[0].Name, mockItems[0].Price, mockItems[0].TaxCode, mockItems[0].Type, mockItems[0].Refundable, mockItems[0].Tax, mockItems[0].Amount).
		AddRow(mockItems[1].ID, mockItems[1].Name, mockItems[1].Price, mockItems[1].TaxCode, mockItems[1].Type, mockItems[1].Refundable, mockItems[1].Tax, mockItems[1].Amount)

	query := `select id, name, price, tax_code,
(case when id = 1 then 'Food & Beverage' when id = 2 then 'Tobacco'  when id = 3 then 'Entertainment' end) as type,
	(case when id = 1 then 'Refundable' when id = 2 then 'Not Refundable'  when id = 3 then 'Not Refundable' end) as refundable,
	(case when id = 1 then round(price * 0.1)
	when id = 2 then round(10 + (0.02 * price))
	when id = 3 then
	case when price BETWEEN 0 AND 100 then 0
	when price >= 100 then round(0.1 * (price - 100))
	end
	end) as tax,

	(case when id = 1 then round((price * 0.1) + price)
	when id = 2 then round((10 + (0.02 * price)) + price)
	when id = 3 then
	case when price BETWEEN 0 AND 100 then round(0 + price)
	when price >= 100 then round((0.1 * (price - 100)) + price )
	end
	end) as amount
	from item ORDER BY id desc `

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := itemRepo.NewMysqlItemRepository(db)

	num := int64(2)
	list, err := a.Fetch(context.TODO(), num)

	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestCalculate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"price_subtotal", "tax_subtotal", "grand_total"}).
		AddRow(2150, 135, 2285)

	query := `select sum(price) as price_subtotal, 
	sum(case when id = 1 then round(price * 0.1)
	when id = 2 then round(10 + (0.02 * price))
	when id = 3 then
	case when price BETWEEN 0 AND 100 then 0
	when price >= 100 then round(0.1 * (price - 100))
	end
	end)  as tax_subtotal,
	sum(case when id = 1 then round((price * 0.1) + price)
	when id = 2 then round((10 + (0.02 * price)) + price)
	when id = 3 then
	case when price BETWEEN 0 AND 100 then round(0 + price)
	when price >= 100 then round((0.1 * (price - 100)) + price )
	end
	end) as grand_total from item `

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := itemRepo.NewMysqlItemRepository(db)

	num := int64(5)
	anItem, err := a.Calculate(context.TODO(), num)
	assert.NoError(t, err)
	assert.NotNil(t, anItem)
}

func TestStore(t *testing.T) {
	ar := &models.Item{
		Name: "Movie",
		Price:  150,
		TaxCode: 3,
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := "INSERT item SET name=\\? , price=\\? , tax_code=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.Name, ar.Price, ar.TaxCode).WillReturnResult(sqlmock.NewResult(3, 1))

	a := itemRepo.NewMysqlItemRepository(db)

	err = a.Store(context.TODO(), ar)
	assert.NoError(t, err)
	assert.Equal(t, int64(12), ar.ID)
}
