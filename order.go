package gokounta

//Order defines an sale from Kounta
type Order struct {
	ID       int     `json:"id"`
	SaleDate string  `json:"created_at"`
	Status   string  `json:"status"`
	Total    float64 `json:"total"`

	Items []OrderLine `json:"lines"`
}

type OrderLine struct {
	TotalTax float64 `json:"line_total_tax"`
}

// GetTotalTax will return the total tax for an order
func (order *Order) GetTotalTax() float64 {
	t := 0.00
	for _, item := range order.Items {
		t += item.TotalTax
	}
	return t
}
