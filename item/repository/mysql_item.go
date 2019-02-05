package repository

import (
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"

	"github.com/shop/item"
	"github.com/shop/models"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type mysqlItemRepository struct {
	Conn *sql.DB
}

// NewMysqlArticleRepository will create an object that represent the article.Repository interface
func NewMysqlItemRepository(Conn *sql.DB) item.Repository {

	return &mysqlItemRepository{Conn}
}

func (m *mysqlItemRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Item, error) {

	rows, err := m.Conn.QueryContext(ctx, query, args...)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer rows.Close()
	result := make([]*models.Item, 0)
	for rows.Next() {
		s := new(models.Item)

		err = rows.Scan(
			&s.ID,
			&s.Name,
			&s.Price,
			&s.TaxCode,
			&s.Type,
			&s.Refundable,
			&s.Tax,
			&s.Amount,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		result = append(result, s)
	}

	return result, nil
}

func (m *mysqlItemRepository) fetchCalculate(ctx context.Context, query string, args ...interface{}) ([]*models.Calculate, error) {

	rows, err := m.Conn.QueryContext(ctx, query, args...)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer rows.Close()
	result := make([]*models.Calculate, 0)
	for rows.Next() {
		s := new(models.Calculate)

		err = rows.Scan(
			&s.PriceSubtotal,
			&s.TaxSubtotal,
			&s.GrandTotal,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		result = append(result, s)
	}

	return result, nil
}

func (m *mysqlItemRepository) Fetch(ctx context.Context, num int64) ([]*models.Item, error) {
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

	res, err := m.fetch(ctx, query)

	if err != nil {
		return nil, err
	}

	return res, err

}

func (m *mysqlItemRepository) Calculate(ctx context.Context, num int64) (*models.Calculate, error){
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

	list, err := m.fetchCalculate(ctx, query)
	if err != nil {
		return nil, err
	}

	a := &models.Calculate{}
	if len(list) > 0 {
		a = list[0]
	} else {
		return nil, models.ErrNotFound
	}

	return a, nil

}


func (m *mysqlItemRepository) Store(ctx context.Context, a *models.Item) error {

	query := `INSERT item SET name=? , price=? , tax_code=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {

		return err
	}

	res, err := stmt.ExecContext(ctx, a.Name, a.Price, a.TaxCode)
	if err != nil {
		return err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return err
	}
	a.ID = lastId
	return nil
}

