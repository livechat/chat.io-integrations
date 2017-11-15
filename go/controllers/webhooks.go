package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type WebhookController struct{}

type WebhookResponse struct {
	ID        string          `json:"webhook_id,omitempty"`
	SecretKey string          `json:"secret_key,omitempty"`
	Action    string          `json:"action,omitempty"`
	Data      interface{}     `json:"-"`
	RawData   json.RawMessage `json:"data,omitempty"`
}

type IncomingEvent struct {
	ChatID string
}

func NewWebhookController() *WebhookController {
	return &WebhookController{}
}

func (w *WebhookController) Handle(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}

	response := &WebhookResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}

	switch response.Action {
	case "incoming_event":
		w.handleIncomingEvent(response.RawData)
	}
}

func (w *WebhookController) handleIncomingEvent(body json.RawMessage) {
	fmt.Println(string(body))
}
