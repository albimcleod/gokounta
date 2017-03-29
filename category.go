package gokounta

//Category is the struct for a Kounta category
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

//Categories is the struct for a list of Category
type Categories []Category
