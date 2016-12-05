package gokounta

// TokenResponse is the response for requesting a token
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

//WebHookRequest is the request structs for creating a webhook
type WebHookRequest struct {
	Topic   string `json:"topic"`
	Address string `json:"address"`
	Format  string `json:"format"`
}
