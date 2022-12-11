package model

type WavFileDetails struct {
	Name           string
	FileSize       int
	Duration       *float64
	NumChannels    *int
	SampleRate     *int
	AudioFormat    *uint16
	AvgBytesPerSec *uint32
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
