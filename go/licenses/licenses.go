package licenses

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"livechat/integration/go/config"
	"livechat/integration/go/customers"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	chat_io "livechat/integration/go/chat.io"
)

type License struct {
	ID           uint64
	AccessToken  string
	RefreshToken string

	Customers *customers.Customers
}

type Licenses struct {
	list       map[uint64]*License
	config     *config.Configuration
	mu         *sync.RWMutex
	httpClient *http.Client
}

type StartChatRequest struct {
	RoutingStatus *RoutingStatus `json:"routing_status"`
}

type RoutingStatus struct {
	Type string `json:"type"`
}

func NewLicenses(config *config.Configuration) *Licenses {
	return &Licenses{
		list:   make(map[uint64]*License, 0),
		config: config,
		mu:     &sync.RWMutex{},
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (l *Licenses) Setup(id uint64, token, refreshToken string) {
	if l.License(id) == nil {
		l.Add(id, token, refreshToken)
	}

	// if err := l.registerWebhooksForLicense(id); err != nil {
	// 	fmt.Println("ERR", err)
	// }

	license := l.License(id)

	// create example customer
	c := license.Customers.Create(id)
	if c != nil {
		log.Println("NEW CUSTOMER!", c.ID, c.AccessToken)

		payload := &StartChatRequest{
			RoutingStatus: &RoutingStatus{
				Type: "license",
			},
		}

		u := &url.URL{}
		qs := u.Query()
		qs.Add("license_id", fmt.Sprintf("%d", id))
		u.RawQuery = qs.Encode()

		// start sample chat
		err := chat_io.CustomerAPI().REST().Send(l.config.Services.External.CustomerAPI.URL, "0.3", "start_chat", c.AccessToken, u, payload)
		if err != nil {
			log.Print(err)
		}
	}
}

func (l *Licenses) Add(id uint64, token, refreshToken string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.list[id] = &License{
		ID:           id,
		AccessToken:  token,
		RefreshToken: refreshToken,
		Customers:    customers.NewCustomers(l.config, token),
	}
}

func (l *Licenses) License(id uint64) *License {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.list[id]
}

type RegisterWebhookRequest struct {
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
	Action      string `json:"action,omitempty"`
	SecretKey   string `json:"secret_key,omitempty"`
}

func (l *Licenses) registerWebhooksForLicense(id uint64) error {
	webhookReq := &RegisterWebhookRequest{
		URL:         "http://webhook.lc4labs.ultrahook.com",
		Description: "Test webhook",
		Action:      "incoming_event",
		SecretKey:   fmt.Sprintf("%d", id),
	}

	raw, err := json.Marshal(webhookReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", l.config.Services.External.ConfigurationAPI.URL+"/v0.2/webhooks/register_webhook", bytes.NewBuffer(raw))
	if err != nil {
		return err
	}

	license := l.License(id)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", license.AccessToken))

	resp, err := l.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode, fmt.Sprintf("Bearer %s", license.AccessToken))

	if resp.StatusCode != http.StatusOK {
		return errors.New("internal_error")
	}

	return nil
}
