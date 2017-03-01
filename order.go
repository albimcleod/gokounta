package gokounta

//Order defines an sale from Kounta
type Order struct {
	ID       int     `json:"id"`
	SaleDate string  `json:"updated_at"`
	Status   string  `json:"status"`
	Total    float64 `json:"total"`

	Items []OrderLine `json:"lines"`
}

//OrderLine defines an line of an order from Kounta
type OrderLine struct {
	TotalTax  float64          `json:"line_total_tax"`
	Product   OrderLineProduct `json:"product"`
	UnitPrice float64          `json:"unit_price"`
	UnitTax   float64          `json:"unit_tax"`
	Quantity  float64          `json:"quantity"`
}

//OrderLineProduct defines an product with an order from Kounta
type OrderLineProduct struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GetTotalTax will return the total tax for an order
func (order *Order) GetTotalTax() float64 {
	t := 0.00
	for _, item := range order.Items {
		t += item.TotalTax
	}
	return t
}
