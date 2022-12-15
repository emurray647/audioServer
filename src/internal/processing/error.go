package processing

import (
	"errors"
)

// predefined errors that we should handle
var (
	invalidFileFormat error = errors.New("invalid file format")
	unknownFileType   error = errors.New("unknown file extension")
	fileAlreadyExists error = errors.New("file already exists")
	fileDoesNotExist  error = errors.New("file does not exist")
)
