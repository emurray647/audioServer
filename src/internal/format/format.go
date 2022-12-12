package format

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/emurray647/audioServer/internal/model"
)

func ParseFile(name string, buffer []byte) (*model.WavFileDetails, error) {

	reader := bytes.NewReader(buffer)

	var result *model.WavFileDetails
	var err error

	fileExtension := filepath.Ext(name)
	switch fileExtension {
	case ".wav":
		result, err = parseWav(reader)
		if err != nil {
			return nil, fmt.Errorf("failed parsing file: %w", err)
		}
	case ".mp3":
		result, err = parseMp3(reader)
		if err != nil {
			return nil, fmt.Errorf("failed parsing file: %w", err)
		}
	default:
		fmt.Println("unknown extension")

	}

	result.Name = name
	result.FileSize = len(buffer)
	return result, nil

}
