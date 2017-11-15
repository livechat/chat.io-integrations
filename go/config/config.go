package config

import (
	"encoding/json"
	"os"
)

type (
	Configuration struct {
		Name    string `json:"name"`
		Port    int    `json:"port"`
		Env     string `json:"env"`
		LogPath string `json:"log_path"`

		Application struct {
			ClientID     string `json:"client_id"`
			ClientSecret string `json:"secret"`
			RedirectURI  string `json:"redirect_uri"`
		} `json:"application"`

		Services struct {
			External struct {
				AgentSSO struct {
					URL string `json:"url"`
				} `json:"agent_sso"`
				CustomerSSO struct {
					URL string `json:"url"`
				} `json:"customer_sso"`
				ConfigurationAPI struct {
					URL string `json:"url"`
				} `json:"configuration_api"`
			} `json:"external"`
		} `json:"services"`
	}
)

func NewConfiguration(fileName string) *Configuration {
	config := &Configuration{}
	config.init(fileName)

	return config
}

func (self *Configuration) init(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&self)
}
