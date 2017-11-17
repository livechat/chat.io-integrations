package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	chat_io "livechat/integration/go/customer_api_chat_example/chat.io"
	chat_io_capi "livechat/integration/go/customer_api_chat_example/chat.io/customer-api"
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

	// parse protocol message
	msg := chat_io.CustomerAPI().ProtocolMessages().IncomingEvent(body)
	if msg == nil {
		return
	}

	// parse event
	event, err := chat_io.CustomerAPI().Events().Unmarshal(msg.RawEvent)
	fmt.Println(event, string(msg.RawEvent), err)
	if err != nil {
		return
	}

	switch chat_io.CustomerAPI().Events().Type(event) {
	case "message":
		message := event.(*chat_io_capi.EventMessage)
		w.handleAction(message.Text)
	default:
		// unsupported type
		return
	}
}

func (w *WebhookController) handleAction(action string) {
	switch action {
	case "hello":
		fmt.Println("Hello there!")
	}
}
