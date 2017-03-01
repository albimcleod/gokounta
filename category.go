package gokounta

//KountaCategory is the struct for a KountaCategory company
type KountaCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type KountaCategories []KountaCategory
