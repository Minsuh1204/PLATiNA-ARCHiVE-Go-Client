package platinaarchivegoclient

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const baseURL = "https://www.platina-archive.app"

func FetchArchive(b64APIKey string) ([]Archive, *APIError) {
	url := baseURL + "/api/v2/get_archive"
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatalf("Error making a new request: %v", err)
	}
	req.Header.Add("X-API-Key", b64APIKey)
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error doing request: %v", err)
	}
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}
		var apiError APIError
		err = json.Unmarshal(bodyBytes, &apiError)
		if err != nil {
			log.Fatalf("Error parsing JSON: %v", err)
		}
		return nil, &apiError
	}
	var archives []Archive
	err = json.Unmarshal(bodyBytes, &archives)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	return archives, nil
}

func FetchClientVersion() ClientVersion {
	url := baseURL + "/api/v1/client_version"
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatalf("API request failed with status code %d. Response: %s", res.StatusCode, bodyBytes)
	}
	var version ClientVersion
	err = json.Unmarshal(bodyBytes, &version)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	return version
}

func FetchPatterns(cache *Cache) ([]Pattern, bool) {
	url := baseURL + "/api/v1/platina_patterns"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error making a new request: %v", err)
	}
	req.Header.Add("If-Modified-Since", cache.PatternsLastModified)
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error doing request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotModified {
		return cache.Patterns, false
	}
	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		log.Fatalf("API request failed with status code %d. Response: %s", res.StatusCode, bodyBytes)
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	var patterns []Pattern
	err = json.Unmarshal(bodyBytes, &patterns)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	return patterns, true
}

func FetchSongs(cache *Cache) ([]Song, bool) {
	url := baseURL + "/api/v1/platina_songs"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error making a new request: %v", err)
	}
	req.Header.Add("If-Modified-Since", cache.SongsLastModified)
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error doing request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotModified {
		return cache.Songs, false
	}
	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		log.Fatalf("API request failed with status code %d. Response: %s", res.StatusCode, bodyBytes)
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	var songs []Song
	err = json.Unmarshal(bodyBytes, &songs)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	return songs, true
}

func Login(name string, password string) (*LoginResult, *APIError) {
	url := baseURL + "/api/v1/login"
	data := map[string]string{"name": name, "password": password}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error composing JSOM: %v", err)
	}
	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error doing request: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}
		var apiError APIError
		err = json.Unmarshal(bodyBytes, &apiError)
		if err != nil {
			log.Fatalf("Error parsing JSON: %v", err)
		}
		return nil, &apiError
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	var result LoginResult
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	return &result, nil
}

func Register(name string, password string) (*RegisterResult, *APIError) {
	url := baseURL + "/api/v1/register"
	data := map[string]string{"name": name, "password": password}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error composing JSON: %v", err)
	}
	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error doing request: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}
		var apiError APIError
		err = json.Unmarshal(bodyBytes, &apiError)
		if err != nil {
			log.Fatalf("Error parsing JSON: %v", err)
		}
		return nil, &apiError
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	var info RegisterResult
	err = json.Unmarshal(bodyBytes, &info)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	return &info, nil
}

func loadCache(cachePath string) (Cache, error) {
	cacheFilePath := filepath.Join(getCacheDirectory(), "cache.json")
	file, err := os.Open(cacheFilePath)
	if err != nil {
		log.Fatalf("Error opening cache file: %v", err)
	}
	defer file.Close()
	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading the file: %v", err)
	}
	var cache Cache
	err = json.Unmarshal(byteValue, &cache)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	return cache, nil
}

func getCacheDirectory() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Error loading user config folder: %v", err)
	}
	return filepath.Join(configDir, "PLATiNA-ARCHiVE")
}
