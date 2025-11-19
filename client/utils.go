package client

import (
	"time"

	"github.com/zalando/go-keyring"
)

const SERVICE_NAME = "PlatinaArchiveClient"
const USERNAME = "main_api_key"

func FormatCurrentTime() string {
	return time.Now().Format("15:04:05")
}

func LoadAPIKey() string {
	key, err := keyring.Get(SERVICE_NAME, USERNAME)
	if err != nil {
		return "" // If no API key is found, return an empty string
	}
	return key
}

func SaveAPIKey(apiKey string) error {
	return keyring.Set(SERVICE_NAME, USERNAME, apiKey)
}
