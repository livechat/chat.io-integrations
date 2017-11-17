package customers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"livechat/integration/go/customer_api_chat_example/config"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type Customers struct {
	list       map[string]*Customer
	mu         *sync.RWMutex
	httpClient *http.Client
	config     *config.Configuration

	accessToken string // integration access_token
}

type Customer struct {
	LicenseID   uint64
	ID          string
	Key         string
	AccessToken string
}

type CustomerIdentityResponse struct {
	Code        string `json:"code,omitempty"`
	CustomerID  string `json:"customer_id,omitempty"`
	CustomerKey string `json:"customer_key,omitempty"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresIn    uint64 `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	LicenseID    uint64 `json:"license_id,omitempty"`
}

func NewCustomers(cfg *config.Configuration, accessToken string) *Customers {
	return &Customers{
		accessToken: accessToken,
		config:      cfg,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Customers) Create(licenseID uint64) *Customer {
	// 1. Request for customer_id, customer_key and code
	identity, err := c.generateChatIOIdentity()
	if err != nil {
		log.Print("3", err)
		return nil
	}

	// 2. Exchange code for access_token to use chat.io Customer API
	accessToken, err := c.exchangeCodeForToken(identity.Code)
	if err != nil {
		log.Print("4", err)
		return nil
	}

	return &Customer{
		LicenseID:   licenseID,
		ID:          identity.CustomerID,
		Key:         identity.CustomerKey,
		AccessToken: accessToken.AccessToken,
	}
}

func (c *Customers) Add(id string, cmr *Customer) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.list[id] = cmr
}

func (c *Customers) Customer(id string) *Customer {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.list[id]
}

func (c *Customers) generateChatIOIdentity() (*CustomerIdentityResponse, error) {

	data := url.Values{}
	data.Add("client_id", c.config.Application.ClientID)
	data.Add("redirect_uri", c.config.Application.RedirectURI)
	data.Add("response_type", "code")

	req, err := http.NewRequest("POST", c.config.Services.External.CustomerSSO.URL+"/", bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println(c.config.Services.External.CustomerSSO.URL+"/", data, c.accessToken)

	if resp == nil {
		return nil, errors.New("internal_error")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("internal_error")
	}

	tokenRes := &CustomerIdentityResponse{}
	if err := json.Unmarshal(body, tokenRes); err != nil {
		return nil, err
	}

	return tokenRes, nil
}

func (c *Customers) exchangeCodeForToken(code string) (*TokenResponse, error) {

	data := url.Values{}
	data.Set("code", code)
	data.Add("client_id", c.config.Application.ClientID)
	data.Add("client_secret", c.config.Application.ClientSecret)
	data.Add("redirect_uri", c.config.Application.RedirectURI)
	data.Add("grant_type", "authorization_code")

	req, err := http.NewRequest("POST", c.config.Services.External.CustomerSSO.URL+"/token", bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp == nil {
		return nil, errors.New("internal_error")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("internal_error")
	}

	tokenRes := &TokenResponse{}
	if err := json.Unmarshal(body, tokenRes); err != nil {
		return nil, err
	}

	return tokenRes, nil
}
