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

type Config struct {
	Reference        Reference        `yaml:"reference" json:"reference"`
	SpeedWidgetPHash string           `yaml:"speedWidgetPHash" json:"speedWidgetPHash"`
	Configs          []ROIConfig      `yaml:"configs" json:"configs"`
	DifficultyColors DifficultyColors `yaml:"difficultyColors" json:"difficultyColors"`
	ColorTolerance   int              `yaml:"colorTolerance" json:"colorTolerance"`
}

type Reference struct {
	Width  int `yaml:"width" json:"width"`
	Height int `yaml:"height" json:"height"`
}

type ROIConfig struct {
	ScreenSize string          `yaml:"screenSize" json:"screenSize"`
	Select     SelectROIConfig `yaml:"select" json:"select"`
	Result     ResultROIConfig `yaml:"result" json:"result"`
}

type SelectROIConfig struct {
	SpeedWidget []int `yaml:"speedWidget" json:"speedWidget"`
	Jacket      []int `yaml:"jacket" json:"jacket"`
	MajorJudge  []int `yaml:"majorJudge" json:"majorJudge"`
	MinorJudge  []int `yaml:"minorJudge" json:"minorJudge"`
	Line        []int `yaml:"line" json:"line"`
	MajorPatch  []int `yaml:"majorPatch" json:"majorPatch"`
	MinorPatch  []int `yaml:"minorPatch" json:"minorPatch"`
	Score       []int `yaml:"score" json:"score"`
	FullCombo   []int `yaml:"fullCombo" json:"fullCombo"`
	MaxPatch    []int `yaml:"maxPatch" json:"maxPatch"`
	Rank        []int `yaml:"rank" json:"rank"`
}

type ResultROIConfig struct {
	Jacket      []int `yaml:"jacket" json:"jacket"`
	Judge       []int `yaml:"judge" json:"judge"`
	Line        []int `yaml:"line" json:"line"`
	Level       []int `yaml:"level" json:"level"`
	Patch       []int `yaml:"patch" json:"patch"`
	Score       []int `yaml:"score" json:"score"`
	Rank        []int `yaml:"rank" json:"rank"`
	NotesArea   []int `yaml:"notesArea" json:"notesArea"`
	TotalNotes  []int `yaml:"totalNotes" json:"totalNotes"`
	PerfectHigh []int `yaml:"perfectHigh" json:"perfectHigh"`
	Perfect     []int `yaml:"perfect" json:"perfect"`
	Great       []int `yaml:"great" json:"great"`
	Good        []int `yaml:"good" json:"good"`
	Miss        []int `yaml:"miss" json:"miss"`
	Difficulty  []int `yaml:"difficulty" json:"difficulty"`
}

type DifficultyColors struct {
	Easy []int `yaml:"easy" json:"easy"`
	Hard []int `yaml:"hard" json:"hard"`
	Over []int `yaml:"over" json:"over"`
	Plus []int `yaml:"plus" json:"plus"`
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
