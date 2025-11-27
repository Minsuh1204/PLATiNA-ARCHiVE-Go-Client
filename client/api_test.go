package client

import (
	"encoding/base64"
	"errors"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func TestFetchArchive(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	password := os.Getenv("ARCHIVE_PASSWORD")
	result, _ := Login("Endeavy", password)
	base64APIKey := base64.StdEncoding.EncodeToString([]byte(result.APIKey))
	archives, apiErr := FetchArchive(base64APIKey)
	if len(archives) != 298 {
		t.Errorf("FetchArchive function does not return every archives: %v", archives)
	}
	if apiErr != nil {
		t.Errorf("FetchArchive funtion return error: %v", apiErr)
	}
}

func TestFetchArchiveInvalidAPIKey(t *testing.T) {
	archives, apiErr := FetchArchive("invalidAPIKey")
	expectedError := APIError{"API key is not encoded correctly"}
	if archives != nil {
		t.Errorf("FetchArchive function does not return nil when API key is invalid: %v", archives)
	}
	if apiErr == nil {
		t.Error("FetchArchive function does not return error when API key is invalid")
	} else {
		var ae *APIError
		if errors.As(apiErr, &ae) {
			if expectedError.Message != ae.Message {
				t.Errorf("FetchArchive function does not return expected error: %v", ae)
			}
		} else {
			t.Errorf("FetchArchive returned unexpected error type: %v", apiErr)
		}
	}
}

func TestFetchClientVersion(t *testing.T) {
	expected := Version{0, 3, 4}
	actual, err := FetchClientVersion()
	if err != nil {
		t.Errorf("FetchClientVersion returned error: %v", err)
	}
	if expected != actual {
		t.Errorf("Version doesn't match, expected: %v, actual: %v", expected, actual)
	}
}

func TestFetchPatternsNoNeedsUpdate(t *testing.T) {
	testPattern := Pattern{0, 4, "EASY", 20, "#Endeavy"}
	cache := Cache{"2025-11-14", "2025-11-06", []Song{}, []Pattern{testPattern}}
	patterns, isUpdated, err := FetchPatterns(&cache)
	if err != nil {
		t.Errorf("FetchPatterns returned error: %v", err)
	}
	if !reflect.DeepEqual(cache.Patterns, patterns) {
		t.Errorf("FetchPatterns function did not return cached patterns: %v", patterns)
	}
	if isUpdated {
		t.Error("FetchPatterns function returned true for update status")
	}
}

func TestFetchPatternsNeedsUpdate(t *testing.T) {
	cache := Cache{"2025-11-01", "2025-11-01", []Song{}, []Pattern{}}
	patterns, isUpdated, err := FetchPatterns(&cache)
	if err != nil {
		t.Errorf("FetchPatterns returned error: %v", err)
	}
	if len(patterns) < 1069 { // Number of songs in DB
		t.Errorf("FetchPatterns function did not give updated patterns list: %v", patterns)
	}
	if !isUpdated {
		t.Error("FetchPatterns function returned false for update status")
	}
}

func TestFetchSongsNoNeedsUpdate(t *testing.T) {
	testSong := Song{0, "example", "artist", "120", "someDLC", "pHash", "plusPHash"}
	cache := Cache{"2025-11-14", "2025-11-14", []Song{testSong}, []Pattern{}}
	songs, isUpdated, err := FetchSongs(&cache)
	if err != nil {
		t.Errorf("FetchSongs returned error: %v", err)
	}
	if !reflect.DeepEqual(cache.Songs, songs) {
		t.Errorf("FetchSongs function did not return cached songs: %v", songs)
	}
	if isUpdated {
		t.Error("FetchSongs function returned true for update status")
	}
}

func TestFetchSongsNeedsUpdate(t *testing.T) {
	cache := Cache{"2025-11-01", "2025-11-14", []Song{}, []Pattern{}}
	songs, isUpdated, err := FetchSongs(&cache)
	if err != nil {
		t.Errorf("FetchSongs returned error: %v", err)
	}
	if len(songs) < 111 { // Number of songs in DB
		t.Errorf("FetchSongs function did not give updated songs list: %v", songs)
	}
	if !isUpdated {
		t.Error("FetchSongs function returned false for update status")
	}
}

func TestRegisterNameAlreadyUsed(t *testing.T) {
	expectedError := APIError{"Name already taken"}
	result, err := Register("Endeavy", "password")
	if err == nil {
		t.Error("Register function does not return error when name is taken")
	} else {
		var ae *APIError
		if errors.As(err, &ae) {
			if expectedError.Message != ae.Message {
				t.Errorf("Register function does not return expected error: %v", ae)
			}
		} else {
			t.Errorf("Register returned unexpected error type: %v", err)
		}
	}
	if result != nil {
		t.Errorf("Register functinon does not return nil when name is taken: %v", *result)
	}
}

func TestRegisterSuccess(t *testing.T) {
	result, err := Register("테스트", "test")
	if err != nil {
		t.Errorf("Register function return error when register is successful: %v", err)
	}
	if result.Name != "테스트" {
		t.Errorf("Register function return wrong name: %v", &result.Name)
	}
	if !strings.HasPrefix(result.APIKey, "테스트::") {
		t.Errorf("Register function returned wrong API key: %v", result.APIKey)
	}
}

func TestLoginSuccess(t *testing.T) {
	result, err := Login("테스트", "test")
	if err != nil {
		t.Errorf("Login function return error when login is successful: %v", err)
	}
	if result.Message != "success" {
		t.Errorf("Login function return message not success: %v", result.Message)
	}
	if !strings.HasPrefix(result.APIKey, "테스트::") {
		t.Errorf("Login function returned wrong API key: %v", result.APIKey)
	}
}

func TestLoginFail(t *testing.T) {
	expectedError := APIError{"로그인 실패"}
	result, err := Login("테스트", "wrong")
	if err == nil {
		t.Error("Login function does not return error when login is fail")
	} else {
		var ae *APIError
		if errors.As(err, &ae) {
			if expectedError.Message != ae.Message {
				t.Errorf("Login function does not return expected error: %v", ae)
			}
		} else {
			t.Errorf("Login returned unexpected error type: %v", err)
		}
	}
	if result != nil {
		t.Errorf("Login function does not return nil when login is fail: %v", *result)
	}
}

func TestUpdateArchiveSuccess(t *testing.T) {
	result, _ := Login("테스트", "test")
	base64APIKey := base64.StdEncoding.EncodeToString([]byte(result.APIKey))
	testArchive := Archive{"테스트", 1, 4, "EASY", 7, 99.1, 100, 360, "2025-11-18", true, false}
	isUpdated, err := UpdateArchive(base64APIKey, testArchive)
	if err != nil {
		t.Errorf("UpdateArchive function return error when update is successful: %v", err)
	}
	if !isUpdated {
		t.Error("UpdateArchive function returned false for updated archive")
	}
}

func TestUpdateArchiveInvalidAPIKey(t *testing.T) {
	testArchive := Archive{"테스트", 1, 4, "EASY", 7, 99.1, 100, 360, "2025-11-18", true, false}
	expectedError := APIError{"API key is not encoded correctly"}
	isUpdated, err := UpdateArchive("invalidAPIKey", testArchive)
	if err == nil {
		t.Error("UpdateArchive function does not return error when API key is invalid")
	} else {
		var ae *APIError
		if errors.As(err, &ae) {
			if expectedError.Message != ae.Message {
				t.Errorf("UpdateArchive function does not return expected error: %v", ae)
			}
		} else {
			t.Errorf("UpdateArchive returned unexpected error type: %v", err)
		}
	}
	if isUpdated {
		t.Error("UpdateArchive function returned true for updated archive")
	}
}

func TestUpdateArchiveInvalidSongID(t *testing.T) {
	result, _ := Login("테스트", "test")
	base64APIKey := base64.StdEncoding.EncodeToString([]byte(result.APIKey))
	testArchive := Archive{"테스트", -99, 4, "EASY", 7, 99.1, 100, 360, "2025-11-18", true, false}
	expectedError := APIError{"Unknown song ID"}
	isUpdated, err := UpdateArchive(base64APIKey, testArchive)
	if err == nil {
		t.Error("UpdateArchive function does not return error when song ID is invalid")
	} else {
		var ae *APIError
		if errors.As(err, &ae) {
			if expectedError.Message != ae.Message {
				t.Errorf("UpdateArchive function does not return expected error: %v", ae)
			}
		} else {
			t.Errorf("UpdateArchive returned unexpected error type: %v", err)
		}
	}
	if isUpdated {
		t.Error("UpdateArchive function returned true for updated archive")
	}
}

func TestUpdateArchiveInvalidLevel(t *testing.T) {
	result, _ := Login("테스트", "test")
	base64APIKey := base64.StdEncoding.EncodeToString([]byte(result.APIKey))
	testArchive := Archive{"테스트", 1, 4, "EASY", -99, 99.1, 100, 360, "2025-11-18", true, false}
	expectedError := APIError{"Invalid level value"}
	isUpdated, err := UpdateArchive(base64APIKey, testArchive)
	if err == nil {
		t.Error("UpdateArchive function does not return error when level is invalid")
	} else {
		var ae *APIError
		if errors.As(err, &ae) {
			if expectedError.Message != ae.Message {
				t.Errorf("UpdateArchive function does not return expected error: %v", ae)
			}
		} else {
			t.Errorf("UpdateArchive returned unexpected error type: %v", err)
		}
	}
	if isUpdated {
		t.Error("UpdateArchive function returned true for updated archive")
	}
}

func TestFetchConfigNoNeedsUpdate(t *testing.T) {
	config := Config{Version: "2099-01-01"}
	isUpdated, err := FetchConfig(&config)
	if err != nil {
		t.Errorf("FetchConfig returned error: %v", err)
	}
	if isUpdated {
		t.Error("FetchConfig function returned true for update status")
	}
}

func TestFetchConfigNeedsUpdate(t *testing.T) {
	config := Config{Version: "2000-01-01"}
	isUpdated, err := FetchConfig(&config)
	if err != nil {
		t.Errorf("FetchConfig returned error: %v", err)
	}
	if !isUpdated {
		t.Error("FetchConfig function returned false for update status")
	}
	if config.Version == "2000-01-01" {
		t.Error("FetchConfig function did not update the config")
	}
}
