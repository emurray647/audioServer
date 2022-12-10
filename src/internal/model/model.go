package model

type WavFileDetails struct {
	Name     string
	Duration int
}

type WavFilesDetailsSlice []WavFileDetails

type WavFile struct {
	Name     string
	Duration int
	URI      string
}
