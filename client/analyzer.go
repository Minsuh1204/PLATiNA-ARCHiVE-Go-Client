package client

import (
	"fmt"
	"image"
	"math"
	"math/big"
	"strconv"

	"github.com/corona10/goimagehash"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

type ScreenType int

const (
	SelectScreen ScreenType = iota
	ResultScreen
)

const pHashThreshold = 2

func AnalyzeScreenshot(cache *Cache, config *Config) (AnalysisReport, error) {
	img, err := LoadImageFromClipboard()
	if err != nil {
		return AnalysisReport{}, fmt.Errorf("failed to load image from clipboard: %v", err)
	}

	// Determine which ROIConfig to use based on screen size

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	screenSize := ScreenSize{Width: width, Height: height}

	var targetConfig ROIConfig
	doScale := false
	found := false
	for _, cfg := range config.Configs {
		if cfg.ScreenSize == screenSize.String() {
			targetConfig = cfg
			found = true
			break
		}
	}

	if !found {
		targetConfig = config.Configs[0] // fallback to default config (1920x1080, reference)
		doScale = true
	}

	screenType, err := decideScreenType(img, targetConfig, doScale, screenSize, config)
	if err != nil {
		return AnalysisReport{}, fmt.Errorf("failed to decide screen type: %v", err)
	}
	jacketMap := buildJacketMap(cache)

	// Extract jacket image
	var baseJacketCoords []int
	switch screenType {
	case SelectScreen:
		baseJacketCoords = targetConfig.Select.Jacket
	case ResultScreen:
		baseJacketCoords = targetConfig.Result.Jacket
	}
	jacketRect := image.Rect(baseJacketCoords[0], baseJacketCoords[1], baseJacketCoords[2], baseJacketCoords[3])
	if doScale {
		x1, y1 := scaleCoordinate(baseJacketCoords[0], baseJacketCoords[1], config.Reference, screenSize)
		x2, y2 := scaleCoordinate(baseJacketCoords[2], baseJacketCoords[3], config.Reference, screenSize)
		jacketRect = image.Rect(x1, y1, x2, y2)
	}
	jacketImage, err := cropImage(img, jacketRect)
	if err != nil {
		return AnalysisReport{}, fmt.Errorf("failed to crop jacket image: %v", err)
	}
	// Calculate pHash of jacket image
	jacketHash, err := goimagehash.PerceptionHash(jacketImage)
	if err != nil {
		return AnalysisReport{}, fmt.Errorf("failed to calculate pHash of jacket image: %v", err)
	}

	// Find best match song
	bestMatchSong, distance, err := bestMatchSong(jacketMap, *jacketHash)
	if err != nil {
		return AnalysisReport{}, fmt.Errorf("failed to find best match song: %v", err)
	}
	// Check if distance is within threshold
	if distance > pHashThreshold {
		return AnalysisReport{}, fmt.Errorf("best match song distance is too high: %d (hash: %v)", distance, jacketHash.GetHash())
	}

	return AnalysisReport{SongObject: bestMatchSong, JacketImage: jacketImage}, nil
}

func decideScreenType(img image.Image, targetConfig ROIConfig, doScale bool, screenSize ScreenSize, config *Config) (ScreenType, error) {
	coords := targetConfig.Select.SpeedWidget
	if len(coords) != 4 {
		return -1, fmt.Errorf("invalid speed widget coordinates")
	}

	rect := image.Rect(coords[0], coords[1], coords[2], coords[3])
	if doScale {
		x1, y1 := scaleCoordinate(coords[0], coords[1], config.Reference, screenSize)
		x2, y2 := scaleCoordinate(coords[2], coords[3], config.Reference, screenSize)
		rect = image.Rect(x1, y1, x2, y2)
	}
	cropped, err := cropImage(img, rect)
	if err != nil {
		return -1, fmt.Errorf("failed to crop image: %v", err)
	}

	knownHashUInt := convertPythonHashToGoHash(config.SpeedWidgetPHash)
	knownHash := goimagehash.NewImageHash(knownHashUInt, goimagehash.PHash)

	calculatedHash, err := goimagehash.PerceptionHash(cropped)
	if err != nil {
		return -1, fmt.Errorf("failed to calculate pHash: %v", err)
	}

	distance, err := calculatedHash.Distance(knownHash)
	if err != nil {
		return -1, fmt.Errorf("failed to calculate hamming distance: %v", err)
	}

	if distance <= pHashThreshold {
		return SelectScreen, nil
	}

	return ResultScreen, nil
}

// Convert Python hash string (hex) to Go hash (uint64)
func convertPythonHashToGoHash(pHash string) uint64 {
	i, ok := new(big.Int).SetString(pHash, 16)
	if !ok {
		return 0
	}
	return i.Uint64()
}

func cropImage(img image.Image, rect image.Rectangle) (image.Image, error) {
	subImager, ok := img.(SubImager)
	if !ok {
		return nil, fmt.Errorf("image does not support cropping")
	}
	return subImager.SubImage(rect), nil
}

func scaleCoordinate(x int, y int, reference ScreenSize, userScreenSize ScreenSize) (int, int) {
	return int(math.Round(float64(x) / float64(reference.Width) * float64(userScreenSize.Width))), int(math.Round(float64(y) / float64(reference.Height) * float64(userScreenSize.Height)))
}

func buildJacketMap(cache *Cache) map[string]Song {
	jacketMap := make(map[string]Song)
	for _, song := range cache.Songs {
		jacketMap[strconv.FormatUint(convertPythonHashToGoHash(song.PHash), 10)] = song
		jacketMap[strconv.FormatUint(convertPythonHashToGoHash(song.PlusPHash), 10)] = song
	}

	return jacketMap
}

func bestMatchSong(jacketMap map[string]Song, hash goimagehash.ImageHash) (Song, int, error) {
	bestMatch := Song{}
	bestDistance := math.MaxInt64
	for compareHash, song := range jacketMap {
		distance, err := hash.Distance(goimagehash.NewImageHash(convertPythonHashToGoHash(compareHash), goimagehash.PHash))
		if err != nil {
			return Song{}, 0, err
		}
		if distance < bestDistance {
			bestDistance = distance
			bestMatch = song
		}
	}
	return bestMatch, bestDistance, nil
}
