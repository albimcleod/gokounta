package gokounta

//KountaProduct is the struct for a KountaCategory company
type KountaProduct struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

//KountaProducts is a slice of KountaProduct
type KountaProducts []KountaProduct
