package gokounta

//WebHook is the request structs for creating a webhook
type WebHook struct {
	ID      string  `json:"id,omitempty"`
	Topic   string  `json:"topic"`
	Address string  `json:"address"`
	Format  string  `json:"format"`
	Filter  Filters `json:"filter,omitempty"`
}

//Filters for the Kount WebHook
type Filters struct {
	SiteID []int `json:"site_ids"`
}

//WebHooks is the struct for a list of WebHook
type WebHooks []WebHook
