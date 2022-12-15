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

// ParseFile parses an audio file and returns an AudioFileDetails object
// name - the name of the file
// buffer - the files contents
func ParseFile(name string, buffer []byte) (*model.AudioFileDetails, error) {
	reader := bytes.NewReader(buffer)

	result := model.AudioFileDetails{
		Name:     name,
		FileSize: len(buffer),
		Format:   "unknown",
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
		parser = unknownParser{}
		err = parser.Parse(&result, reader)
		if err != nil {
			return nil, UnknownFormat
		}
	}

	return &result, nil
}

// Generic interface to implement for all the different methods of parsing audio files
type Parser interface {
	Parse(*model.AudioFileDetails, io.ReadSeeker) error
}

// A parser implementation for if we do not know the format ahead of time
type unknownParser struct {
}

func (p unknownParser) Parse(details *model.AudioFileDetails, reader io.ReadSeeker) error {
	var parser Parser
	var err error

	parser = WavParser{}
	err = parser.Parse(details, reader)
	if err == nil {
		// this parsed correctly, so must be a wav
		return nil
	}

	parser = Mp3Parser{}
	err = parser.Parse(details, reader)
	if err == nil {
		// this parsed correctly, so must be a mp3
		return nil
	}

	return UnknownFormat
}
