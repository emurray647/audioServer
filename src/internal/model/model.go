package model

// Metadata associated with an audio file
type AudioFileDetails struct {
	Name           string   `json:"name"`
	FileSize       int      `json:"file_size"` // in bytes
	Format         string   `json:"format"`    // "wav" or "mp3"
	Duration       *float64 `json:"duration"`  // in seconds
	NumChannels    *int     `json:"num_channels"`
	SampleRate     *int     `json:"sample_rate"` // samples per second
	AvgBytesPerSec *uint32  `json:"avg_bytes_per_second"`
}

// Convenience typedef for a set of AudioFileDetails
type AudioFileDetailsSlice []AudioFileDetails

// All data associated with an AudioFile
type AudioFile struct {
	AudioFileDetails
	// The location where the actual style is being stored
	URI string
}

// Status to provide to the caller about the result of their call
type StatusMessage struct {
	StatusCode int    `json:"status_code"`
	Success    bool   `json:"success"`
	Message    string `json:"message"`
}
