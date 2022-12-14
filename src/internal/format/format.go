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
		Format:   "unknown",
	}
	var err error

	var parser Parser
	fileExtension := filepath.Ext(name)
	fmt.Println(fileExtension)
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
		parser = unknownParser{}
		err = parser.Parse(&result, reader)
		if err != nil {
			return nil, UnknownFormat
		}
	}

	return &result, nil

}

type Parser interface {
	Parse(*model.WavFileDetails, io.ReadSeeker) error
}

type unknownParser struct {
}

func (p unknownParser) Parse(details *model.WavFileDetails, reader io.ReadSeeker) error {
	var parser Parser
	var err error

	fmt.Println("wav")

	parser = WavParser{}
	err = parser.Parse(details, reader)
	if err == nil {
		// this parsed correctly, so must be a wav
		return nil
	}

	fmt.Println("mp3")
	parser = Mp3Parser{}
	err = parser.Parse(details, reader)
	if err == nil {
		// this parsed correctly, so must be a mp3
		return nil
	}

	return UnknownFormat
}
