package platinaarchivegoclient

type APIError struct {
	Message string `json:"msg"`
}

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

type Cache struct {
	SongsLastModified    string    `json:"songsLastModified"`
	PatternsLastModified string    `json:"patternsLastModified"`
	Songs                []Song    `json:"songs"`
	Patterns             []Pattern `json:"patterns"`
}

type ClientVersion struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
}

type LoginResult struct {
	Message string `json:"msg"`
	APIKey  string `json:"key"`
}

type Pattern struct {
	SongID     int    `json:"songID"`
	Line       int    `json:"line"`
	Difficulty string `json:"difficulty"`
	Level      int    `json:"level"`
	Designer   string `json:"designer"`
}

type RegisterResult struct {
	Name   string `json:"name"`
	APIKey string `json:"key"`
}

type Song struct {
	ID        int    `json:"songID"`
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	BPM       string `json:"BPM"`
	DLC       string `json:"DLC"`
	PHash     string `json:"pHash"`
	PlusPHash string `json:"plusPHash"`
}
