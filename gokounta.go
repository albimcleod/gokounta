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
	baseURL               = "https://api.kounta.com"
	webHookURL            = "v1/companies/%v/webhooks"
	tokenURL              = "v1/token.json"
	companiesURL          = "v1/companies/me"
	sitesURL              = "v1/companies/%v/sites"
	webHookTopicSale      = "orders/completed"
	categoriesURL         = "v1/companies/%v/categories"
	categoriesProductsURL = "/v1/companies/%v/categories/%v/products"
	ordersURL             = "v1/companies/%v/sites/%v/orders/pending.json"
	staffURL              = "v1/companies/%v/staff"
	ordersCompleteURL     = "v1/companies/%v/sites/%v/orders/complete.json"
	ordersSingleURL       = "v1/companies/%v/orders/%v.json"
)

var (
	defaultSendTimeout = time.Second * 30
)

// Kounta The main struct of this package
type Kounta struct {
	StoreCode    string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Timeout      time.Duration
}

// NewClient will create a Kounta client with default values
func NewClient(code string, clientID string, clientSecret string, redirectURL string) *Kounta {
	return &Kounta{
		StoreCode:    code,
		Timeout:      defaultSendTimeout,
		RedirectURL:  redirectURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

// AccessToken will get a new access token
func (v *Kounta) AccessToken() (string, string, error) {

	data := url.Values{}
	data.Set("code", v.StoreCode)
	data.Add("client_secret", v.ClientSecret)
	data.Add("client_id", v.ClientID)
	data.Add("redirect_uri", v.RedirectURL)
	data.Add("grant_type", "authorization_code")

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = tokenURL
	urlStr := fmt.Sprintf("%v", u)

	fmt.Printf("AccessToken %v %v\n", urlStr, data)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, _ := client.Do(r)
	fmt.Println(res.Status)

	rawResBody, err := ioutil.ReadAll(res.Body)

	fmt.Printf("AccessToken Body %v \n", string(rawResBody))

	if err != nil {
		return "", "", fmt.Errorf("%v", string(rawResBody))
	}

	if res.StatusCode == 200 {
		resp := &TokenResponse{}
		if err := json.Unmarshal(rawResBody, resp); err != nil {
			return "", "", err
		}

		return resp.AccessToken, resp.RefreshToken, nil
	}

	return "", "", fmt.Errorf("Failed to get refresh token: %s", res.Status)
}

// RefreshToken will get a new refresh token
func (v *Kounta) RefreshToken(refreshtoken string) (string, string, error) {

	data := url.Values{}
	data.Set("refresh_token", refreshtoken)
	data.Add("client_id", v.ClientID)
	data.Add("client_secret", v.ClientSecret)
	data.Add("grant_type", "refresh_token")
	data.Add("redirect_uri", v.RedirectURL)

	u, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return "", "", err
	}

	u.Path = tokenURL
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	r, err := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", "", err
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

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
		if err := json.Unmarshal(rawResBody, resp); err != nil {
			return "", "", err
		}

		return resp.AccessToken, resp.RefreshToken, nil
	}

	return "", "", fmt.Errorf("Error requesting access token")
}

// GetCompany will return the authenticated company
func (v *Kounta) GetCompany(token string) (*Company, error) {
	client := &http.Client{}
	client.CheckRedirect = checkRedirectFunc

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = companiesURL
	urlStr := fmt.Sprintf("%v", u)

	r, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	r.Header = http.Header(make(map[string][]string))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Authorization", "Bearer "+token)

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	//	fmt.Println("GetCompany Body", string(rawResBody))

	if res.StatusCode == 200 {
		var resp Company
		err = json.Unmarshal(rawResBody, &resp)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	}
	return nil, fmt.Errorf("Failed to get Kounta Company %s", res.Status)

}

// GetSites will return the sites of the authenticated company
//not finished
func (v *Kounta) GetSites(token string, company string) (Sites, error) {
	client := &http.Client{}
	client.CheckRedirect = checkRedirectFunc

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = fmt.Sprintf(sitesURL, company)
	urlStr := fmt.Sprintf("%v", u)

	r, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	r.Header = http.Header(make(map[string][]string))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Authorization", "Bearer "+token)

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("GetSites Body", string(rawResBody))

	if res.StatusCode == 200 {
		var resp Sites

		err = json.Unmarshal(rawResBody, &resp)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	return nil, fmt.Errorf("Failed to get Kounta Company %s", res.Status)

}

// GetStaff will return the staff of the authenticated company
func (v *Kounta) GetStaff(token string, company string) (Staffs, error) {
	client := &http.Client{}
	client.CheckRedirect = checkRedirectFunc

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = fmt.Sprintf(staffURL, company)
	urlStr := fmt.Sprintf("%v", u)

	r, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	r.Header = http.Header(make(map[string][]string))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Authorization", "Bearer "+token)

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("GetStaff Body", string(rawResBody))

	if res.StatusCode == 200 {
		var resp Staffs

		err = json.Unmarshal(rawResBody, &resp)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	return nil, fmt.Errorf("Failed to get Kounta Staff %s", res.Status)

}

// GetWebHooks will return the webhooks of the authenticated company
func (v *Kounta) GetWebHooks(token string, company string) (WebHooks, error) {
	client := &http.Client{}
	client.CheckRedirect = checkRedirectFunc

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = fmt.Sprintf(webHookURL+".json", company)
	urlStr := fmt.Sprintf("%v", u)

	r, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	r.Header = http.Header(make(map[string][]string))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Authorization", "Bearer "+token)

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 200 {
		var resp WebHooks

		fmt.Println(string(rawResBody))

		err = json.Unmarshal(rawResBody, &resp)

		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	return nil, fmt.Errorf("Failed to get Kounta Web Hooks %s", res.Status)

}

// CreateSaleWebHook will init the sales hook for the Kounta store
func (v *Kounta) CreateSaleWebHook(token string, company string, webhook WebHook) error {

	fmt.Println("CreateSaleWebHook", token, company, webhook)

	b, err := json.Marshal(webhook)
	if err != nil {
		return err
	}

	u, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return err
	}

	u.Path = fmt.Sprintf(webHookURL+".json", company)
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	client.CheckRedirect = checkRedirectFunc

	r, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	r.Header.Add("Authorization", "Bearer "+token)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Content-Length", strconv.Itoa(len(b)))

	res, err := client.Do(r)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("Failed init sale webhooks %s", res.Status)
	}

	return nil
}

// DeleteSaleWebHook will init the sales hook for the Kounta store
func (v *Kounta) DeleteSaleWebHook(token string, company string, id int) error {

	fmt.Println("UpdateSaleWebHook", token, company, id)

	u, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return err
	}

	u.Path = fmt.Sprintf(webHookURL+"/"+strconv.Itoa(id)+".json", company)
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	client.CheckRedirect = checkRedirectFunc

	r, err := http.NewRequest("DELETE", urlStr, nil) //, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	r.Header.Add("Authorization", "Bearer "+token)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Content-Length", "0") // strconv.Itoa(len(b)))

	res, err := client.Do(r)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("Failed init sale webhooks %s", res.Status)
	}

	return nil
}

// GetCategories will return the categories of the authenticated company
func (v *Kounta) GetCategories(token string, company string) (Categories, error) {
	client := &http.Client{}
	client.CheckRedirect = checkRedirectFunc

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = fmt.Sprintf(categoriesURL, company)
	urlStr := fmt.Sprintf("%v", u)

	r, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	r.Header = http.Header(make(map[string][]string))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Authorization", "Bearer "+token)

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 200 {
		var resp []Category

		err = json.Unmarshal(rawResBody, &resp)

		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	return nil, fmt.Errorf("Failed to get Kounta Categories %s", res.Status)

}

// GetProducts will return the products of the authenticated company
func (v *Kounta) GetProducts(token string, company string, categoryID string) (KountaProducts, error) {

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = fmt.Sprintf(categoriesProductsURL, company, categoryID)

	urlStr := fmt.Sprintf("%v", u)

	results := new(KountaProducts)

	for urlStr != "" {

		resp := new(KountaProducts)

		*resp, _, urlStr = v.callProduct(urlStr, token)

		*results = append(*results, *resp...)

		fmt.Println("X-Next-Page", urlStr, len(*results))
	}

	return *results, nil
}

func (v *Kounta) callProduct(urlStr string, token string) (KountaProducts, error, string) {
	client := &http.Client{}
	client.CheckRedirect = checkRedirectFunc

	r, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err, ""
	}

	r.Header = http.Header(make(map[string][]string))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Authorization", "Bearer "+token)

	res, err := client.Do(r)
	if err != nil {
		return nil, err, ""
	}

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err, ""
	}

	if res.StatusCode == 200 {
		var resp []KountaProduct

		err = json.Unmarshal(rawResBody, &resp)

		if err != nil {
			return nil, err, ""
		}
		return resp, nil, res.Header.Get("X-Next-Page")
	}
	return nil, fmt.Errorf("Failed to get Kounta Products %s", res.Status), ""
}

// GetOrders will return the orders of the authenticated company
func (v *Kounta) GetOrders(token string, company string, siteID string) ([]Order, error) {
	client := &http.Client{}
	client.CheckRedirect = checkRedirectFunc

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = fmt.Sprintf(ordersURL, company, siteID)
	urlStr := fmt.Sprintf("%v", u)

	r, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	r.Header = http.Header(make(map[string][]string))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Authorization", "Bearer "+token)

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 200 {
		var resp []Order

		//fmt.Println(string(rawResBody))

		err = json.Unmarshal(rawResBody, &resp)

		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	return nil, fmt.Errorf("Failed to get Kounta Categories %s", res.Status)

}

// GetOrders will return the orders of the authenticated company
func (v *Kounta) GetOrdersComplete(token string, company string, siteID string, start string) ([]Order, error) {
	client := &http.Client{}
	client.CheckRedirect = checkRedirectFunc

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = fmt.Sprintf(ordersCompleteURL, company, siteID)
	urlStr := fmt.Sprintf("%v", u)

	//urlStr += "?created_gte=2018-08-28"
	if start != "" {
		urlStr += "?start=" + start
	}

	//fmt.Println(urlStr)

	r, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	r.Header = http.Header(make(map[string][]string))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Authorization", "Bearer "+token)

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 200 {
		var resp []Order

		fmt.Println(res.Header["X-Next-Page"])

		//	fmt.Println(string(rawResBody))

		err = json.Unmarshal(rawResBody, &resp)

		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	fmt.Println(string(rawResBody))
	return nil, fmt.Errorf("Failed to get Kounta Categories %s", res.Status)

}

// GetOrders will return the orders of the authenticated company
func (v *Kounta) GetOrdersSingle(token string, company string, orderID string) (*Order, error) {
	client := &http.Client{}
	client.CheckRedirect = checkRedirectFunc

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = fmt.Sprintf(ordersSingleURL, company, orderID)
	urlStr := fmt.Sprintf("%v", u)

	fmt.Println(urlStr)

	r, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	r.Header = http.Header(make(map[string][]string))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Authorization", "Bearer "+token)

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 200 {
		resp := Order{}

		//	fmt.Println(string(rawResBody))

		err = json.Unmarshal(rawResBody, &resp)

		if err != nil {
			return nil, err
		}
		return &resp, nil
	}
	fmt.Println(string(rawResBody))
	return nil, fmt.Errorf("Failed to get Kounta Sale %s", res.Status)

}

func checkRedirectFunc(req *http.Request, via []*http.Request) error {
	if req.Header.Get("Authorization") == "" {
		req.Header.Add("Authorization", via[0].Header.Get("Authorization"))
	}
	return nil
}
