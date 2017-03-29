package gokounta

//Site is the struct for a Kounta Site
type Site struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

//Sites is the struct for a list of Site
type Sites []Site
