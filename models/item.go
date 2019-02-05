package models

type Item struct {
	ID        	int64     	`json:"id"`
	Name      	string    	`json:"name" validate:"required"`
	Price     	int    		`json:"price" validate:"required"`
	TaxCode   	int    		`json:"tax_code" validate:"required"`
	Type 		string		`json:"type"`
	Refundable 	string		`json:"refundable"`
	Tax			int			`json:"tax"`
	Amount		int			`json:"amount"`
}
