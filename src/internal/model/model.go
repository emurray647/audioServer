package model

type WavFileDetails struct {
	Name           string   `json:"name"`
	FileSize       int      `json:"file_size"`
	Duration       *float64 `json:"duration"`
	NumChannels    *int     `json:"num_channels"`
	SampleRate     *int     `json:"sample_rate"`
	AudioFormat    *uint16  `json:"audio_format"`
	AvgBytesPerSec *uint32  `json:"avg_bytes_per_second"`
}

type WavFilesDetailsSlice []WavFileDetails

type WavFile struct {
	// Name     string
	// Duration int
	WavFileDetails
	URI string
}

type StatusMessage struct {
	StatusCode int    `json:"status_code"`
	Success    bool   `json:"success"`
	Message    string `json:"message"`
}
