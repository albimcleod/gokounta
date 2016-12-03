package gokounta

import (
	"bytes"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	AUTHORIZE_URL        = "https://my.kounta.com/authorize"
	BASE_URL             = "https://api.kounta.com/v1/"
	CONTENT_TYPE         = "application/json"
	CONTENT_TYPE_TOKEN   = "application/x-www-form-urlencoded"
	POST                 = "POST"
	HEADER_CONTENT_TYPE  = "Content-Type"
	HEADER_AUTHORIZATION = "Authorization"
	URL_TOKEN            = "token.json"
)

var (
	DEFAULT_SEND_TIMEOUT time.Duration = time.Second * 30
)

// Kounta The main struct of this package
type Kounta struct {
	StoreCode    string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Timeout      time.Duration
}

//
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type KountaWebHookResult struct {
	Type       string `json:"type"`
	URL        string `json:"url"`
	Active     bool   `json:"active"`
	RetailerID string `json:"retailer_id"`
	ID         string `json:"id"`
}

type KountaWebHookRequest struct {
	Type   string `json:"type"`
	URL    string `json:"url"`
	Active bool   `json:"active"`
}

// NewClient will create a Kounta client with default values
func NewClient(code string, clientID string, clientSecret string, redirectURL string) *Kounta {
	return &Kounta{
		StoreCode:    code,
		Timeout:      DEFAULT_SEND_TIMEOUT,
		RedirectURL:  redirectURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

// AccessToken .. wil get a new access token
func (v *Kounta) AccessToken() (string, string, error) {

	data := url.Values{}
	data.Set("code", v.StoreCode)
	data.Add("client_secret", v.ClientSecret)
	data.Add("client_id", v.ClientID)
	data.Add("response_type", "code")
	data.Add("redirect_uri", v.RedirectURL)

	u, _ := url.ParseRequestURI(AUTHORIZE_URL)
	//	u.Path = v.TokenURL
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, _ := client.Do(r)
	fmt.Println(res.Status)

	rawResBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", "", err
	}

	if res.StatusCode >= 400 {
		return "", "", fmt.Errorf("Failed to get refresh token: %s", res.Status)
	}

	if res.StatusCode == 200 {
		resp := &TokenResponse{}
		err = json.Unmarshal(rawResBody, resp)
		if err != nil {
			return "", "", err
		}

		return resp.AccessToken, resp.RefreshToken, nil
	}

	return "", "", fmt.Errorf("Error requesting access token")
}

// RefreshToken .. wil get a new fresh token
func (v *Kounta) RefreshToken(refreshtoken string) (string, string, error) {

	data := url.Values{}
	data.Set("code", v.StoreCode)
	data.Set("refresh_token", refreshtoken)
	data.Add("client_id", v.ClientID)
	data.Add("client_secret", v.ClientSecret)
	data.Add("grant_type", "refresh_token")
	data.Add("redirect_uri", v.RedirectURL)

	u, _ := url.ParseRequestURI(BASE_URL)
	u.Path = URL_TOKEN
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	fmt.Println("urlStr", urlStr, data)

	res, _ := client.Do(r)
	fmt.Println(res.Status)

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", "", err
	}

	fmt.Println("BODY", string(rawResBody))

	if res.StatusCode >= 400 {
		return "", "", fmt.Errorf("Failed to get refresh token: %s", res.Status)
	}

	if res.StatusCode == 200 {
		resp := &TokenResponse{}
		err = json.Unmarshal(rawResBody, resp)
		if err != nil {
			return "", "", err
		}

		return resp.AccessToken, resp.RefreshToken, nil
	}

	return "", "", fmt.Errorf("Error requesting access token")
}

/*
// InitSaleWebHook will init the sales hook for the Kounta store
func (v *Kounta) InitSaleWebHook(token string, uri string) error {

	webhook := KountaWebHookRequest{
		Type:   WEB_HOOK_SALE,
		URL:    uri,
		Active: true,
	}

	b, err := json.Marshal(webhook)
	if err != nil {
		return err
	}

	data := url.Values{}
	data.Set("data", string(b))

	u, _ := url.ParseRequestURI(v.BaseURL)
	u.Path = v.WebHookURL
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
	r.Header.Add("Authorization", "Bearer "+token)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, _ := client.Do(r)
	fmt.Println(res.Status)

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("Failed init sale webhooks %s", res.Status)
	}

	if res.StatusCode == 200 {
		resp := &TokenResponse{}
		err = json.Unmarshal(rawResBody, resp)
		if err != nil {
			return err
		}
	}

	return nil
}

// RevokeExistingWebHooks will init the sales hook for the Kounta store
func (v *Kounta) RevokeExistingWebHooks(token string, storeID string) error {
	u, _ := url.ParseRequestURI(v.BaseURL)
	u.Path = v.WebHookURL
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	r, _ := http.NewRequest("GET", urlStr, nil)

	r.Header.Add("Authorization", "Bearer "+token)

	res, _ := client.Do(r)
	fmt.Println(res.Status)

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("Failed revoke sale webhooks %s", res.Status)
	}

	if res.StatusCode == 200 {
		var resp []KountaWebHookResult
		err = json.Unmarshal(rawResBody, &resp)
		if err != nil {
			return err
		}

		for _, webhook := range resp {
			if strings.Contains(webhook.URL, storeID) {
				err = v.RevokeWebHook(token, webhook)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// RevokeWebHook will init the sales hook for the Kounta store
func (v *Kounta) RevokeWebHook(token string, webhook KountaWebHookResult) error {

	u, _ := url.ParseRequestURI(v.BaseURL)
	u.Path = v.WebHookURL + "/" + webhook.ID
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	r, _ := http.NewRequest("DELETE", urlStr, nil)

	r.Header.Add("Authorization", "Bearer "+token)

	res, _ := client.Do(r)
	fmt.Println(res.Status)

	_, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("Failed revoke sale webhook %s", res.Status)
	}
	return nil
}

// GetSales will init the sales hook for the Kounta store
func (v *Kounta) GetSales(token string, storeID string) error {
	u, _ := url.ParseRequestURI(v.BaseURL)
	u.Path = v.SalesURL
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}

	urlStr = urlStr + "?since=2016-11-28T20:00:15.000Z&page=1"
	r, _ := http.NewRequest("GET", urlStr, nil)

	r.Header.Add("Authorization", "Bearer "+token)

	res, _ := client.Do(r)
	fmt.Println(res.Status)

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("Failed get sales %s", res.Status)
	}

	if res.StatusCode == 200 {
		substring := rawResBody[0:200]

		fmt.Println("BODY", string(substring))

		var resp KountaSaleResult
		err = json.Unmarshal(rawResBody, &resp)
		if err != nil {
			fmt.Println(err)
			return err
		}

		//	total := 0.00

		shr, smin := 9, 30
		ehr, emin := 12, 30

		for _, sale := range resp.Sales {

			timezone, _ := time.LoadLocation("Australia/Adelaide")

			ld := sale.GetSaleDate().In(timezone)

			if ld.Hour() > shr || (ld.Hour() == shr && ld.Minute() > smin) {
				if ld.Hour() < ehr || (ld.Hour() == ehr && ld.Minute() <= emin) {

					fmt.Println(ld, sale, ld.Hour(), ld.Minute())
					//	total += sale.Totals.Total + sale.Totals.TotalTax
				}
			}
		}

		//	fmt.Println("total=", total)
	}

	return nil
}
*/
