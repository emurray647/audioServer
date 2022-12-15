package format

import (
	"io"

	"github.com/emurray647/audioServer/internal/model"
	"github.com/go-audio/wav"
)

// A parser to parse wav files
type WavParser struct{}

func (WavParser) Parse(details *model.AudioFileDetails, reader io.ReadSeeker) error {
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
	details.AvgBytesPerSec = &decoder.AvgBytesPerSec

	return nil
}
