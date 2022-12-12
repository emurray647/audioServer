package format

import (
	"io"

	"github.com/emurray647/audioServer/internal/model"
	"github.com/tcolgate/mp3"
)

func parseMp3(reader io.ReadSeeker) (*model.WavFileDetails, error) {

	decoder := mp3.NewDecoder(reader)

	duration := 0.0

	var frame mp3.Frame
	skipped := 0

	for {
		if err := decoder.Decode(&frame, &skipped); err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		duration += float64(frame.Duration().Seconds())

		// fmt.Println(frame.Samples())
		// fmt.Println(float64(frame.Samples()) / float64(frame.Duration().Seconds()))
	}

	result := model.WavFileDetails{
		Duration: &duration,
		// AvgBytesPerSec: ,
	}

	return &result, nil
}
