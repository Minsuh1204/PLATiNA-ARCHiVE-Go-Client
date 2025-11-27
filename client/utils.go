package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"time"

	"github.com/zalando/go-keyring"
	"golang.design/x/clipboard"
	"gopkg.in/yaml.v3"
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
	defaultCache := Cache{
		SongsLastModified:    "2025-04-11",
		PatternsLastModified: "2025-04-11",
	}

	if !fileExists(cacheFilePath) {
		err := updateCache(&defaultCache)
		return defaultCache, err
	}
	file, err := os.OpenFile(cacheFilePath, os.O_RDWR, 0644)
	if err != nil {
		updateCache(&defaultCache)
		return defaultCache, fmt.Errorf("error opening cache file: %w", err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&defaultCache); err != nil {
		updateCache(&defaultCache)
		return defaultCache, fmt.Errorf("error parsing JSON: %v", err)
	}
	return defaultCache, nil
}

func updateCache(cache *Cache) error {
	songs, songsLoaded, err := FetchSongs(cache)
	if !songsLoaded {
		return fmt.Errorf("failed to fetch songs from server")
	}
	if err != nil {
		return fmt.Errorf("error fetching songs from server: %v", err)
	}
	patterns, patternsLoaded, err := FetchPatterns(cache)
	if !patternsLoaded {
		return fmt.Errorf("failed to fetch patterns from server")
	}
	if err != nil {
		return fmt.Errorf("error fetching patterns from server: %v", err)
	}
	cache.Songs = songs
	cache.Patterns = patterns
	cache.SongsLastModified = time.Now().Format(time.RFC3339)
	cache.PatternsLastModified = time.Now().Format(time.RFC3339)
	err = saveCache(cache)
	if err != nil {
		return fmt.Errorf("error saving cache: %v", err)
	}
	return nil
}

func saveCache(cache *Cache) error {
	os.MkdirAll(filepath.Join(getCacheDirectory(), "cache"), os.ModeDir)
	cacheFilePath := filepath.Join(getCacheDirectory(), "cache", "db.json")
	file, err := os.OpenFile(cacheFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening cache file: %v", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(cache); err != nil {
		return fmt.Errorf("error writing JSON: %v", err)
	}
	return nil
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

func LoadConfig() (Config, error) {
	// Creates the cache directory if it doesn't exist
	os.MkdirAll(getCacheDirectory(), os.ModeDir)
	// Creates the empty cache file if it doesn't exist
	cacheFilePath := filepath.Join(getCacheDirectory(), "config.yaml")

	defaultConfig := Config{Version: "2025-04-11"}

	if !fileExists(cacheFilePath) {
		res, err := FetchConfig(&defaultConfig)
		if !res {
			return Config{}, fmt.Errorf("error fetching config: %v", err)
		}
		SaveConfig(&defaultConfig)
		return defaultConfig, nil
	}
	file, err := os.OpenFile(cacheFilePath, os.O_RDONLY, 0644)
	if err != nil {
		FetchConfig(&defaultConfig)
		SaveConfig(&defaultConfig)
		return defaultConfig, fmt.Errorf("error opening cache file: %v", err)
	}
	defer file.Close()

	var config Config
	err = yaml.NewDecoder(file).Decode(&config)
	file.Close()
	if err != nil {
		FetchConfig(&defaultConfig)
		SaveConfig(&defaultConfig)
		return defaultConfig, fmt.Errorf("error parsing YAML: %w", err)
	}
	FetchConfig(&config)
	SaveConfig(&config)
	return config, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func SaveConfig(config *Config) error {
	cacheFilePath := filepath.Join(getCacheDirectory(), "config.yaml")
	file, err := os.OpenFile(cacheFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening cache file: %v", err)
	}
	defer file.Close()

	if err := yaml.NewEncoder(file).Encode(config); err != nil {
		return fmt.Errorf("error writing YAML: %v", err)
	}
	return nil
}

func LoadImageFromClipboard() (image.Image, error) {
	err := clipboard.Init()
	if err != nil {
		return nil, fmt.Errorf("error initializing clipboard: %v", err)
	}
	imgData := clipboard.Read(clipboard.FmtImage)
	if imgData == nil {
		return nil, fmt.Errorf("no image data found in clipboard")
	}

	// Decode the image data
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image from clipboard: %v", err)
	}
	return img, nil
}
