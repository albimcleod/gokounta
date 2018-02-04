package gokounta

//Order defines a sale from Kounta
type Order struct {
	ID             int64         `json:"id"`
	SaleDate       string        `json:"created_at"`
	UpdateDate     string        `json:"updated_at"`
	Status         string        `json:"status"`
	Notes          string        `json:"notes"`
	Total          float64       `json:"total"`
	PriceVariation float64       `json:"price_variation"`
	Customer       OrderCustomer `json:"customer"`
	SiteID         float64       `json:"site_id"`

	Items    []OrderLine    `json:"lines"`
	Payments []OrderPayment `json:"payments"`
}

//OrderCustomer defines  line of an order from Kounta
type OrderCustomer struct {
	ID        int64  `json:"id"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
}

//OrderLine defines a line of an order from Kounta
type OrderLine struct {
	Product        OrderLineProduct `json:"product"`
	UnitPrice      float64          `json:"unit_price"`
	UnitTax        float64          `json:"unit_tax"`
	LineTotal      float64          `json:"line_total_ex_tax"`
	LineTotalTax   float64          `json:"line_total_tax"`
	Quantity       float64          `json:"quantity"`
	PriceVariation float64          `json:"price_variation"`
	Modifiers      []int            `json:"modifiers"`
}

//OrderLineProduct defines a product within an order from Kounta
type OrderLineProduct struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

//OrderPayment defines a payment of an order from Kounta
type OrderPayment struct {
	Number int                `json:"number"`
	Amount float64            `json:"amount"`
	Method OrderPaymentMethod `json:"method"`
}

//OrderPaymentMethod defines a payment method within an order from Kounta
type OrderPaymentMethod struct {
	ID   int64  `json:"id"`
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
