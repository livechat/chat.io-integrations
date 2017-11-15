package main

import (
	"fmt"
	"livechat/integration/go/config"
	"livechat/integration/go/controllers"
	"livechat/integration/go/licenses"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	cfg := config.NewConfiguration("./config/config.json")
	licenses := licenses.NewLicenses(cfg)

	router := mux.NewRouter()
	authController := controllers.NewAuthController(cfg, licenses)
	webhookController := controllers.NewWebhookController()

	router.HandleFunc("/", authController.Auth)
	router.HandleFunc("/webhook", webhookController.Handle)

	fmt.Println(fmt.Sprintf(":%d", cfg.Port))
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)
}
