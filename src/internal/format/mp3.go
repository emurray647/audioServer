package format

import (
	"io"

	"github.com/emurray647/audioServer/internal/model"
	"github.com/tcolgate/mp3"
)

type Mp3Parser struct{}

func (Mp3Parser) Parse(details *model.WavFileDetails, reader io.ReadSeeker) error {
	decoder := mp3.NewDecoder(reader)

	duration := 0.0
	samples := 0

	var frame mp3.Frame
	skipped := 0

	for {
		if err := decoder.Decode(&frame, &skipped); err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		duration += float64(frame.Duration().Seconds())
		samples += frame.Samples()

		// fmt.Println(frame.)
		// frame.Header().

	}

	details.Format = "mp3"
	details.Duration = &duration
	avgBytesPerSec := uint32(float64(details.FileSize) / duration)
	details.AvgBytesPerSec = &avgBytesPerSec

	// TODO: not sure with the mp3 logic here... I think it might be right, but not sure
	numChannels := 1
	channelMode := frame.Header().ChannelMode()
	if channelMode == mp3.DualChannel || channelMode == mp3.JointStereo {
		numChannels = 2
	}
	details.NumChannels = &numChannels

	sampleRate := samples * numChannels / int(duration)
	details.SampleRate = &sampleRate
	// TODO end

	return nil
}
