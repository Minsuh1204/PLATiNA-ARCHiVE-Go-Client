package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = "https://www.platina-archive.app"

// FetchArchive retrieves the user's archive (play history) from the server.
// It requires a base64 encoded API key for authentication.
// Returns a slice of Archive structs or an error if the request fails.
func FetchArchive(b64APIKey string) ([]Archive, error) {
	url := baseURL + "/api/v2/get_archive"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error making new request: %w", err)
	}
	req.Header.Add("X-API-Key", b64APIKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error doing request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var apiError APIError
		if err := json.NewDecoder(res.Body).Decode(&apiError); err != nil {
			return nil, fmt.Errorf("error parsing JSON error response: %w", err)
		}
		return nil, &apiError
	}

	var archives []Archive
	if err := json.NewDecoder(res.Body).Decode(&archives); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}
	return archives, nil
}

// FetchClientVersion retrieves the current client version from the server.
// Returns a ClientVersion struct or an error if the request fails.
func FetchClientVersion() (ClientVersion, error) {
	url := baseURL + "/api/v1/client_version"
	res, err := http.Get(url)
	if err != nil {
		return ClientVersion{}, fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ClientVersion{}, fmt.Errorf("API request failed with status code %d", res.StatusCode)
	}

	var version ClientVersion
	if err := json.NewDecoder(res.Body).Decode(&version); err != nil {
		return ClientVersion{}, fmt.Errorf("error parsing JSON: %w", err)
	}
	return version, nil
}

// FetchPatterns retrieves the list of patterns from the server.
// It uses the provided cache to check for updates using the If-Modified-Since header.
// Returns a slice of Pattern structs, a boolean indicating if the list was updated, or an error.
func FetchPatterns(cache *Cache) ([]Pattern, bool, error) {
	url := baseURL + "/api/v1/platina_patterns"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, false, fmt.Errorf("error making new request: %w", err)
	}
	req.Header.Add("If-Modified-Since", cache.PatternsLastModified)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, false, fmt.Errorf("error doing request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotModified {
		return cache.Patterns, false, nil
	}
	if res.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("API request failed with status code %d", res.StatusCode)
	}

	var patterns []Pattern
	if err := json.NewDecoder(res.Body).Decode(&patterns); err != nil {
		return nil, false, fmt.Errorf("error parsing JSON: %w", err)
	}
	return patterns, true, nil
}

// FetchSongs retrieves the list of songs from the server.
// It uses the provided cache to check for updates using the If-Modified-Since header.
// Returns a slice of Song structs, a boolean indicating if the list was updated, or an error.
func FetchSongs(cache *Cache) ([]Song, bool, error) {
	url := baseURL + "/api/v1/platina_songs"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, false, fmt.Errorf("error making new request: %w", err)
	}
	req.Header.Add("If-Modified-Since", cache.SongsLastModified)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, false, fmt.Errorf("error doing request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotModified {
		return cache.Songs, false, nil
	}
	if res.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("API request failed with status code %d", res.StatusCode)
	}

	var songs []Song
	if err := json.NewDecoder(res.Body).Decode(&songs); err != nil {
		return nil, false, fmt.Errorf("error parsing JSON: %w", err)
	}
	return songs, true, nil
}

// Login authenticates a user with the given name and password.
// Returns a LoginResult struct containing the API key or an error if login fails.
func Login(name string, password string) (*LoginResult, error) {
	url := baseURL + "/api/v1/login"
	data := map[string]string{"name": name, "password": password}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error composing JSON: %w", err)
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error doing request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var apiError APIError
		if err := json.NewDecoder(res.Body).Decode(&apiError); err != nil {
			return nil, fmt.Errorf("error parsing JSON error response: %w", err)
		}
		return nil, &apiError
	}

	var result LoginResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}
	return &result, nil
}

// Register registers a new user with the given name and password.
// Returns a RegisterResult struct containing the API key or an error if registration fails.
func Register(name string, password string) (*RegisterResult, error) {
	url := baseURL + "/api/v1/register"
	data := map[string]string{"name": name, "password": password}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error composing JSON: %w", err)
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error doing request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var apiError APIError
		if err := json.NewDecoder(res.Body).Decode(&apiError); err != nil {
			return nil, fmt.Errorf("error parsing JSON error response: %w", err)
		}
		return nil, &apiError
	}

	var info RegisterResult
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}
	return &info, nil
}

// UpdateArchive updates the user's archive with a new play record.
// It requires a base64 encoded API key for authentication.
// Returns true if the update was successful, or an error if it failed.
func UpdateArchive(b64APIKey string, archive Archive) (bool, error) {
	url := baseURL + "/api/v2/update_archive"
	jsonData, err := json.Marshal(archive)
	if err != nil {
		return false, fmt.Errorf("error composing JSON: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("error making new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-API-Key", b64APIKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error doing request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var apiError APIError
		if err := json.NewDecoder(res.Body).Decode(&apiError); err != nil {
			return false, fmt.Errorf("error parsing JSON error response: %w", err)
		}
		return false, &apiError
	}
	return true, nil
}
