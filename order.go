package gokounta

//Order defines an sale from Kounta
type Order struct {
	ID       int     `json:"id"`
	SaleDate string  `json:"created_at"`
	Status   string  `json:"status"`
	Total    float64 `json:"total"`
	TotalTax float64 `json:"total_tax"`
}
