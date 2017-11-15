package main

import (
	"fmt"
	"livechat/integration/config"
	"livechat/integration/controllers"
	"livechat/integration/licenses"
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

	fmt.Println(fmt.Sprintf(":%d", cfg.Port))
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)
}
