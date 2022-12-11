package processing

import "errors"

var (
	invalidWAVError   error = errors.New("invalid WAV file")
	fileAlreadyExists error = errors.New("file already exists")

	fileDoesNotExist error = errors.New("file does not exist")
)
