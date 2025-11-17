package platinaarchivegoclient

import (
	"encoding/base64"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func TestFetchArchive(t *testing.T) {
	err := godotenv.Load()
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

func TestFetchClientVersion(t *testing.T) {
	expected := ClientVersion{0, 3, 4}
	actual := FetchClientVersion()
	if expected != actual {
		t.Errorf("Version doesn't match, expected: %v, actual: %v", expected, actual)
	}
}

func TestFetchPatternsNoNeedsUpdate(t *testing.T) {
	testPattern := Pattern{0, 4, "EASY", 20, "#Endeavy"}
	cache := Cache{"2025-11-14", "2025-11-06", []Song{}, []Pattern{testPattern}}
	patterns, isUpdated := FetchPatterns(&cache)
	if !reflect.DeepEqual(cache.Patterns, patterns) {
		t.Errorf("FetchPatterns function did not return cached patterns: %v", patterns)
	}
	if isUpdated {
		t.Error("FetchPatterns function returned true for update status")
	}
}

func TestFetchPatternsNeedsUpdate(t *testing.T) {
	cache := Cache{"2025-11-01", "2025-11-01", []Song{}, []Pattern{}}
	patterns, isUpdated := FetchPatterns(&cache)
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
	songs, isUpdated := FetchSongs(&cache)
	if !reflect.DeepEqual(cache.Songs, songs) {
		t.Errorf("FetchSongs function did not return cached songs: %v", songs)
	}
	if isUpdated {
		t.Error("FetchSongs function returned true for update status")
	}
}

func TestFetchSongsNeedsUpdate(t *testing.T) {
	cache := Cache{"2025-11-01", "2025-11-14", []Song{}, []Pattern{}}
	songs, isUpdated := FetchSongs(&cache)
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
	if expectedError != *err {
		t.Errorf("Register function does not return error when name is taken: %v", *err)
	}
	if result != nil {
		t.Errorf("Register functinon does not return nil when name is taken: %v", *result)
	}
}

func TestRegisterSuccess(t *testing.T) {
	result, err := Register("테스트", "test")
	if err != nil {
		t.Errorf("Register function return error when register is successful: %v", *err)
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
		t.Errorf("Login function return error when login is successful: %v", *err)
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
	if expectedError != *err {
		t.Errorf("Login function does not return expected error: %v", *err)
	}
	if result != nil {
		t.Errorf("Login function does not return nil when login is fail: %v", *result)
	}
}
