package models

type Calculate struct {
	PriceSubtotal     	int    		`json:"price_subtotal"`
	TaxSubtotal   		int    		`json:"tax_subtotal"`
	GrandTotal   		int    		`json:"grand_total"`
}
