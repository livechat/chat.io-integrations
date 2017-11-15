package controllers

import "net/http"

type WebhookController struct{}

func NewWebhookController() *WebhookController {
	return &WebhookController{}
}

func (w *WebhookController) Handle(rw http.ResponseWriter, req *http.Request) {
}
