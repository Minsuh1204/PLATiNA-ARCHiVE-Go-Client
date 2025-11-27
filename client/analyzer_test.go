package client

import (
	"image"
	"image/color"
	"os"
	"strconv"
	"testing"
)

// MockSubImager is a mock image that implements SubImager interface
type MockSubImager struct {
	image.RGBA
}

func (m *MockSubImager) SubImage(r image.Rectangle) image.Image {
	return m.RGBA.SubImage(r)
}

func TestCropImage(t *testing.T) {
	// Create a 100x100 image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	// Fill with blue
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			img.Set(x, y, color.RGBA{0, 0, 255, 255})
		}
	}

	// Crop a 10x10 area at (10, 10)
	rect := image.Rect(10, 10, 20, 20)
	cropped, err := cropImage(img, rect)
	if err != nil {
		t.Fatalf("cropImage failed: %v", err)
	}

	if cropped.Bounds() != rect {
		t.Errorf("expected bounds %v, got %v", rect, cropped.Bounds())
	}
}

func TestDecideScreenTypeWithRealImages(t *testing.T) {
	// Load real images
	selectImg, err := loadImage("../testing/select.png")
	if err != nil {
		t.Fatalf("failed to load select.png: %v", err)
	}
	resultImg, err := loadImage("../testing/result.png")
	if err != nil {
		t.Fatalf("failed to load result.png: %v", err)
	}

	// Create config with real values from config.yaml
	config := &Config{
		SpeedWidgetPHash: "c0c73d38273ed2c3",
		Reference: ScreenSize{
			Width:  1920,
			Height: 1080,
		},
		Configs: []ROIConfig{
			{
				ScreenSize: "1920x1080",
				Select: SelectROIConfig{
					SpeedWidget: []int{30, 908, 119, 932},
				},
			},
		},
	}

	// Test Select Screen
	screenType, err := decideScreenType(selectImg, config.Configs[0], false, config.Reference, config)
	if err != nil {
		t.Errorf("decideScreenType failed for select.png: %v", err)
	}
	if screenType != SelectScreen {
		t.Errorf("expected SelectScreen, got %v", screenType)
	}

	// Test Result Screen
	screenType, err = decideScreenType(resultImg, config.Configs[0], false, config.Reference, config)
	if err != nil {
		// It might return error if it's not a select screen (which is expected behavior for now as we only detect SelectScreen)
		// But if it returns an error, we should check if it's the expected "screen not recognized" error or something else.
		// For now, let's assume it should return ResultScreen if we implement it, or error if not.
		// Based on current implementation: returns -1 and error if not SelectScreen.
		// Wait, the user asked to test if it correctly decides ResultScreen too.
		// But current implementation only returns SelectScreen or error/-1.
		// Let's check the implementation of decideScreenType again.
		// Ah, I see I changed it to return ResultScreen in the last turn (Step 326).
		// So it should return ResultScreen.
		t.Errorf("decideScreenType failed for result.png: %v", err)
	}
	if screenType != ResultScreen {
		t.Errorf("expected ResultScreen, got %v", screenType)
	}
}

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	return img, err
}

func TestBuildJacketMap(t *testing.T) {
	// Mock Cache with sample songs
	cache := &Cache{
		Songs: []Song{
			{
				ID:        1,
				Title:     "Test Song 1",
				PHash:     "c0c73d38273ed2c3", // Hex string
				PlusPHash: "a1b2c3d4e5f60708", // Hex string
			},
			{
				ID:        2,
				Title:     "Test Song 2",
				PHash:     "1234567890abcdef",
				PlusPHash: "", // Empty hash
			},
		},
	}

	// Call buildJacketMap
	jacketMap := buildJacketMap(cache)

	// Verify map size
	// Song 1 has 2 valid hashes, Song 2 has 1 valid hash (assuming empty string results in 0 which might be added or ignored depending on implementation, let's check convertPythonHashToGoHash behavior for empty string.
	// convertPythonHashToGoHash("") -> new(big.Int).SetString("", 16) -> error -> returns 0.
	// So empty hash becomes "0".
	// If multiple songs have empty hash, they will overwrite each other at key "0".
	// Let's assume for this test we expect 3 entries if "0" is included.

	// Helper to convert hex to decimal string key
	toKey := func(hex string) string {
		val := convertPythonHashToGoHash(hex)
		return strconv.FormatUint(val, 10)
	}

	// Check Song 1 PHash
	key1 := toKey("c0c73d38273ed2c3")
	if song, exists := jacketMap[key1]; !exists {
		t.Errorf("Expected map to contain key %s for Song 1 PHash", key1)
	} else if song.ID != 1 {
		t.Errorf("Expected Song 1 for key %s, got ID %d", key1, song.ID)
	}

	// Check Song 1 PlusPHash
	key2 := toKey("a1b2c3d4e5f60708")
	if song, exists := jacketMap[key2]; !exists {
		t.Errorf("Expected map to contain key %s for Song 1 PlusPHash", key2)
	} else if song.ID != 1 {
		t.Errorf("Expected Song 1 for key %s, got ID %d", key2, song.ID)
	}

	// Check Song 2 PHash
	key3 := toKey("1234567890abcdef")
	if song, exists := jacketMap[key3]; !exists {
		t.Errorf("Expected map to contain key %s for Song 2 PHash", key3)
	} else if song.ID != 2 {
		t.Errorf("Expected Song 2 for key %s, got ID %d", key3, song.ID)
	}
}
