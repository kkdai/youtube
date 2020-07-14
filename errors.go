package youtube

import (
	"errors"
)

var (
	ErrEmptyStreamList            = errors.New("Empty Stream list")
	ErrItagNotFound               = errors.New("Invalid itag value, please specify correct value.")
	ErrCipherNotFound             = errors.New("cipher not found")
	ErrInvalidCharactersInVideoId = errors.New("invalid characters in video id")
	ErrVideoIdMinLength           = errors.New("the video id must be at least 10 characters long")
)
