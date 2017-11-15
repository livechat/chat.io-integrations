package licenses

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"livechat/integration/config"
	"net/http"
	"sync"
	"time"
)

type License struct {
	ID           uint64
	AccessToken  string
	RefreshToken string
}

type Licenses struct {
	list       map[uint64]*License
	config     *config.Configuration
	mu         *sync.RWMutex
	httpClient *http.Client
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

	if err := l.registerWebhooksForLicense(id); err != nil {
		fmt.Println("ERR", err)
	}
}

func (l *Licenses) Add(id uint64, token, refreshToken string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.list[id] = &License{
		ID:           id,
		AccessToken:  token,
		RefreshToken: refreshToken,
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
