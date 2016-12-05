package gokounta

import (
	"bytes"

	//	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	AUTHORIZE_URL = "https://my.kounta.com/authorize"
	BASE_URL      = "https://api.kounta.com"

	URL_WEB_HOOK = "/v1/companies/%v/webhooks.json"

	URL_TOKEN     = "v1/token.json"
	URL_COMPANIES = "/v1/companies/me.json"

	CONTENT_TYPE         = "application/json"
	CONTENT_TYPE_TOKEN   = "application/x-www-form-urlencoded"
	POST                 = "POST"
	HEADER_CONTENT_TYPE  = "Content-Type"
	HEADER_AUTHORIZATION = "Authorization"

	WEB_HOOK_SALE = "orders/completed"
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

// TokenResponse is the response for requesting a token
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

//KountaWebHookResult is the result structs for creating a webhook
type KountaWebHookResult struct {
	Type       string `json:"type"`
	URL        string `json:"url"`
	Active     bool   `json:"active"`
	RetailerID string `json:"retailer_id"`
	ID         string `json:"id"`
}

//KountaWebHookRequest is the request structs for creating a webhook
type KountaWebHookRequest struct {
	Topic   string `json:"topic"`
	Address string `json:"address"`
	Format  string `json:"format"`
}

//KountaCompany is the struct for a Kounta company
type KountaCompany struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

//KountaSale defines an sale from Kounta
type KountaSale struct {
	ID       string  `json:"id"`
	SaleDate string  `json:"created_at"`
	Status   string  `json:"status"`
	Total    float64 `json:"total"`
	TotalTax float64 `json:"total_tax"`
}

func (obj *KountaSale) GetSaleDate() time.Time {
	d := obj.SaleDate

	if !strings.Contains(d, "T") {
		d = strings.Replace(d, " ", "T", 1)
	}

	if !strings.Contains(d, "Z") {
		d = d + "Z"
	}

	t1, _ := time.Parse(time.RFC3339, d)
	return t1
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
	data.Add("redirect_uri", v.RedirectURL)
	data.Add("grant_type", "authorization_code")

	u, _ := url.ParseRequestURI(BASE_URL)
	u.Path = URL_TOKEN
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, _ := client.Do(r)
	fmt.Println(res.Status)

	rawResBody, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(rawResBody))

	if err != nil {
		return "", "", err
	}

	if res.StatusCode == 200 {
		resp := &TokenResponse{}
		err = json.Unmarshal(rawResBody, resp)
		if err != nil {
			return "", "", err
		}

		return resp.AccessToken, resp.RefreshToken, nil
	}

	return "", "", fmt.Errorf("Failed to get refresh token: %s", res.Status)
}

// RefreshToken .. wil get a new fresh token
func (v *Kounta) RefreshToken(refreshtoken string) (string, string, error) {

	data := url.Values{}
	//	data.Set("code", v.StoreCode)
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

// InitSaleWebHook will init the sales hook for the Kounta store
func (v *Kounta) InitSaleWebHook(token string, company string, uri string) error {

	fmt.Println("InitSaleWebHook", token, company, uri)

	webhook := KountaWebHookRequest{
		Topic:   WEB_HOOK_SALE,
		Address: uri,
		Format:  "json",
	}

	b, err := json.Marshal(webhook)
	if err != nil {
		return err
	}

	//	data := url.Values{}
	//	data.Set("data", string(b))

	hookURL := fmt.Sprintf(URL_WEB_HOOK, company)

	u, _ := url.ParseRequestURI(BASE_URL)
	u.Path = hookURL
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	client.CheckRedirect = checkRedirectFunc

	//	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer(b)) // <-- URL-encoded payload
	r.Header.Add("Authorization", "Bearer "+token)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Content-Length", strconv.Itoa(len(b)))

	res, _ := client.Do(r)
	fmt.Println(res.Status)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println("InitSaleWebHook Body", string(resBody))

	if res.StatusCode >= 400 {
		return fmt.Errorf("Failed init sale webhooks %s", res.Status)
	}

	if res.StatusCode == 200 {

	}

	return nil
}

// GetCompany will return the authenticated company
func (v *Kounta) GetCompany(token string) (*KountaCompany, error) {
	client := &http.Client{}
	client.CheckRedirect = checkRedirectFunc

	r, _ := http.NewRequest("GET", "https://api.kounta.com/v1/companies/me", nil)

	r.Header = http.Header(make(map[string][]string))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Authorization", "Bearer "+token)

	fmt.Println("GetCompany URL=", r.URL)
	fmt.Println("GetCompany TOKEN=", token)
	fmt.Println("GetCompany HEADER=", r.Header)

	res, _ := client.Do(r)
	fmt.Println(res.Status)

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("GetCompany Body", string(rawResBody))

	if res.StatusCode == 200 {
		var resp KountaCompany
		err = json.Unmarshal(rawResBody, &resp)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	}
	return nil, fmt.Errorf("Failed to get Kounta Company %s", res.Status)

}

func checkRedirectFunc(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", via[0].Header.Get("Authorization"))
	return nil
}

/*
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
