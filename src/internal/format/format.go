package format

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/emurray647/audioServer/internal/model"
)

var UnknownFormat error = errors.New("unknown file format")
var InvalidFile error = errors.New("invalid audio file")

func ParseFile(name string, buffer []byte) (*model.WavFileDetails, error) {

	reader := bytes.NewReader(buffer)

	// var result *model.WavFileDetails
	result := model.WavFileDetails{
		Name:     name,
		FileSize: len(buffer),
	}
	var err error

	var parser Parser
	fileExtension := filepath.Ext(name)
	switch fileExtension {
	case ".wav":
		parser = WavParser{}
		err = parser.Parse(&result, reader)
		if err != nil {
			return nil, fmt.Errorf("failed parsing wav file: %w", InvalidFile)
		}
	case ".mp3":
		parser = Mp3Parser{}
		err = parser.Parse(&result, reader)
		if err != nil {
			return nil, fmt.Errorf("failed parsing mp3 file: %w", InvalidFile)
		}
	default:
		return nil, UnknownFormat
	}

	// result.Name = name
	// result.FileSize = len(buffer)
	return &result, nil

}

type Parser interface {
	Parse(*model.WavFileDetails, io.ReadSeeker) error
}
