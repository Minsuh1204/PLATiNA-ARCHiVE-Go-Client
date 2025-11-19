package client

import "image"

// APIError represents an error returned by the API.
type APIError struct {
	Message string `json:"msg"`
}

// Error returns the error message.
func (e *APIError) Error() string {
	return e.Message
}

type AnalysisReport struct {
	SongObject    Song
	PatternObject Pattern
	JacketImage   image.Image
	Judge         float64
	Score         float64
	Patch         float64
	Rank          string
	FullCombo     bool
	MaxPatch      bool
}

// Archive represents a user's play record for a song.
type Archive struct {
	Decoder    string  `json:"decoder"`
	SongID     int     `json:"song_id"`
	Line       int     `json:"line"`
	Difficulty string  `json:"difficulty"`
	Level      int     `json:"level"`
	Judge      float64 `json:"judge"`
	Score      float64 `json:"score"`
	Patch      float64 `json:"patch"`
	DecodedAt  string  `json:"decoded_at"`
	FullCombo  bool    `json:"is_full_combo"`
	MaxPatch   bool    `json:"is_max_patch"`
}

// Cache represents the local cache for songs and patterns.
type Cache struct {
	SongsLastModified    string    `json:"songsLastModified"`
	PatternsLastModified string    `json:"patternsLastModified"`
	Songs                []Song    `json:"songs"`
	Patterns             []Pattern `json:"patterns"`
}

// ClientVersion represents the version of the client.
type ClientVersion struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
}

// LoginResult represents the response from the login API.
type LoginResult struct {
	Message string `json:"msg"`
	APIKey  string `json:"key"`
}

// Pattern represents a chart/pattern for a song.
type Pattern struct {
	SongID     int    `json:"songID"`
	Line       int    `json:"line"`
	Difficulty string `json:"difficulty"`
	Level      int    `json:"level"`
	Designer   string `json:"designer"`
}

// RegisterResult represents the response from the registration API.
type RegisterResult struct {
	Name   string `json:"name"`
	APIKey string `json:"key"`
}

// Song represents a song in the game.
type Song struct {
	ID        int    `json:"songID"`
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	BPM       string `json:"BPM"`
	DLC       string `json:"DLC"`
	PHash     string `json:"pHash"`
	PlusPHash string `json:"plusPHash"`
}
