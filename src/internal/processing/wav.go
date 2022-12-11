package processing

import (
	"bytes"
	"fmt"

	"github.com/emurray647/audioServer/internal/model"
	"github.com/go-audio/wav"
)

type invalidFileError struct{}

func (i *invalidFileError) Error() string {
	return "file is not a valid WAV"
}

func parseWav(name string, buffer []byte) (*model.WavFileDetails, error) {
	// reader := bytes.NewReader(buffer)
	// reader.Seek()
	decoder := wav.NewDecoder(bytes.NewReader(buffer))

	fmt.Println(len(buffer))

	decoder.ReadInfo()

	if !decoder.IsValidFile() {
		// return nil, &invalidFileError{}
		return nil, invalidWAVError
	}

	details := &model.WavFileDetails{
		Name:     name,
		FileSize: len(buffer),
	}

	duration, err := decoder.Duration()
	if err == nil {
		seconds := duration.Seconds()
		details.Duration = &seconds
	}

	details.NumChannels = &decoder.Format().NumChannels
	details.SampleRate = &decoder.Format().SampleRate

	details.AudioFormat = &decoder.WavAudioFormat
	details.AvgBytesPerSec = &decoder.AvgBytesPerSec

	// fmt.Println(decoder.SampleRate)

	// fmt.Println(decoder.IsValidFile())
	// decoder.WavAudioFormat

	// decoder.ReadMetadata()
	// if decoder.Err() != nil {
	// 	fmt.Println(decoder.Err().Error())
	// }
	// fmt.Println(decoder.Metadata.Title)

	return details, nil
}
