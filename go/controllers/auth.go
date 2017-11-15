package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"livechat/integration/go/config"
	"livechat/integration/go/licenses"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	InternalError      error = errors.New("Internal error")
	InvalidAccessToken error = errors.New("Invalid access token")
)

type AuthController struct {
	licenses   *licenses.Licenses
	config     *config.Configuration
	httpClient *http.Client
}

type TokenResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresIn    uint64 `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	LicenseID    uint64 `json:"license_id,omitempty"`
}

func NewAuthController(cfg *config.Configuration, licenses *licenses.Licenses) *AuthController {
	return &AuthController{
		config:   cfg,
		licenses: licenses,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (a *AuthController) Auth(rw http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")

	// 1. Exchange code for token
	tk, err := a.exchangeCodeForToken(code)
	if err != nil {
		fmt.Println(tk, err)
		rw.WriteHeader(500)
		return
	}

	// 2. Get Token info
	token, err := a.getTokenInfo(fmt.Sprintf("%s %s", tk.TokenType, tk.AccessToken))
	if err != nil {
		fmt.Println(token, err)
		rw.WriteHeader(500)
		return
	}

	// 3. Setup license config
	a.licenses.Setup(token.LicenseID, token.AccessToken, token.RefreshToken)
}

func (a *AuthController) exchangeCodeForToken(code string) (*TokenResponse, error) {

	data := url.Values{}
	data.Set("code", code)
	data.Add("client_id", a.config.Application.ClientID)
	data.Add("client_secret", a.config.Application.ClientSecret)
	data.Add("redirect_uri", a.config.Application.RedirectURI)
	data.Add("grant_type", "authorization_code")

	req, err := http.NewRequest("POST", a.config.Services.External.AgentSSO.URL+"/token", bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp == nil {
		return nil, InternalError
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, InternalError
	}

	tokenRes := &TokenResponse{}
	if err := json.Unmarshal(body, tokenRes); err != nil {
		return nil, err
	}

	return tokenRes, nil
}

func (a *AuthController) getTokenInfo(token string) (*TokenResponse, error) {
	req, err := http.NewRequest("GET", a.config.Services.External.AgentSSO.URL+"/info", nil)
	if err != nil {
		return nil, InternalError
	}

	req.Header.Set("Authorization", token)

	resp, err := a.httpClient.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, InvalidAccessToken
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, InternalError
	}

	tokenRes := &TokenResponse{}
	if err := json.Unmarshal(body, tokenRes); err != nil {
		return nil, InternalError
	}

	if resp == nil {
		return nil, InternalError
	}
	if resp.StatusCode != http.StatusOK {
		return nil, InvalidAccessToken
	}

	return tokenRes, nil
}
