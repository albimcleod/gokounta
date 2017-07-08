package gokounta

//Order defines an sale from Kounta
type Order struct {
	ID             int     `json:"id"`
	SaleDate       string  `json:"created_at"`
	UpdateDate     string  `json:"updated_at"`
	Status         string  `json:"status"`
	Total          float64 `json:"total"`
	PriceVariation float64 `json:"price_variation"`

	Items    []OrderLine    `json:"lines"`
	Payments []OrderPayment `json:"payments"`
}

//OrderLine defines an line of an order from Kounta
type OrderLine struct {
	Product        OrderLineProduct `json:"product"`
	UnitPrice      float64          `json:"unit_price"`
	UnitTax        float64          `json:"unit_tax"`
	LineTotal      float64          `json:"line_total_ex_tax"`
	LineTotalTax   float64          `json:"line_total_tax"`
	Quantity       float64          `json:"quantity"`
	PriceVariation float64          `json:"price_variation"`
}

//OrderLineProduct defines an product with an order from Kounta
type OrderLineProduct struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

//OrderPayment defines an line of an order from Kounta
type OrderPayment struct {
	Number int                `json:"number"`
	Amount float64            `json:"amount"`
	Method OrderPaymentMethod `json:"method"`
}

//OrderLineProduct defines an product with an order from Kounta
type OrderPaymentMethod struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetTotalTax will return the total tax for an order
func (order *Order) GetTotalTax() float64 {
	t := 0.00
	for _, item := range order.Items {
		t += item.LineTotalTax
	}
	return t
}

// CalculateDiscount will return the total discount for an order item
func (oi *OrderLine) CalculateDiscount() float64 {
	t := 0.00
	if oi.PriceVariation > 0 && oi.PriceVariation < 1 {
		return ((oi.LineTotal + oi.LineTotalTax) / oi.PriceVariation) - (oi.LineTotal + oi.LineTotalTax)
	}
	return t
}
