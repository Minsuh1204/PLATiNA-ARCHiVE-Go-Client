package client

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

// LoadCache loads the cache from the local file system.
// Returns a Cache struct or an error if loading fails.
func LoadCache(cacheFileName string) (Cache, error) {
	// Creates the cache directory if it doesn't exist
	os.MkdirAll(getCacheDirectory(), os.ModeDir)
	// Creates the empty cache file if it doesn't exist
	cacheFilePath := filepath.Join(getCacheDirectory(), cacheFileName)
	file, err := os.OpenFile(cacheFilePath, os.O_CREATE|os.O_RDONLY|os.O_WRONLY, 0644)
	if err != nil {
		return Cache{}, fmt.Errorf("error opening cache file: %w", err)
	}
	defer file.Close()

	var cache Cache
	if err := json.NewDecoder(file).Decode(&cache); err != nil {
		return Cache{}, fmt.Errorf("error parsing JSON: %w", err)
	}
	return cache, nil
}

// getCacheDirectory returns the directory where the cache is stored.
// It panics if the user config directory cannot be determined.
func getCacheDirectory() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(fmt.Sprintf("error loading user config folder: %v", err))
	}
	return filepath.Join(configDir, "PLATiNA-ARCHiVE")
}
