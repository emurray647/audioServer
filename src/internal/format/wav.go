package format

import (
	"io"

	"github.com/emurray647/audioServer/internal/model"
	"github.com/go-audio/wav"
)

// type InvalidFileError struct{}

// func (i *InvalidFileError) Error() string {
// 	return "file is not a valid WAV"
// }

type WavParser struct{}

func (WavParser) Parse(details *model.WavFileDetails, reader io.ReadSeeker) error {
	decoder := wav.NewDecoder(reader)
	decoder.ReadInfo()

	if !decoder.IsValidFile() {
		return InvalidFile
	}

	duration, err := decoder.Duration()
	if err == nil {
		seconds := duration.Seconds()
		details.Duration = &seconds
	}

	details.Format = "wav"
	details.NumChannels = &decoder.Format().NumChannels
	details.SampleRate = &decoder.Format().SampleRate

	details.AudioFormat = &decoder.WavAudioFormat
	details.AvgBytesPerSec = &decoder.AvgBytesPerSec

	return nil
}

// func parseWav(reader io.ReadSeeker) (*model.WavFileDetails, error) {

// 	// func ParseWav(name string, buffer []byte) (*model.WavFileDetails, error) {

// 	// decoder := wav.NewDecoder(bytes.NewReader(buffer))
// 	decoder := wav.NewDecoder(reader)

// 	// fmt.Println(len(buffer))

// 	decoder.ReadInfo()

// 	if !decoder.IsValidFile() {
// 		return nil, InvalidFile
// 	}

// 	details := &model.WavFileDetails{}

// 	duration, err := decoder.Duration()
// 	if err == nil {
// 		seconds := duration.Seconds()
// 		details.Duration = &seconds
// 	}

// 	details.NumChannels = &decoder.Format().NumChannels
// 	details.SampleRate = &decoder.Format().SampleRate

// 	details.AudioFormat = &decoder.WavAudioFormat
// 	details.AvgBytesPerSec = &decoder.AvgBytesPerSec

// 	return details, nil
// }
