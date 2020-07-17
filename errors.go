package youtube

import (
	"errors"
)

var (
	ErrEmptyStreamList            = errors.New("empty Stream list")
	ErrItagNotFound               = errors.New("invalid itag value, please specify correct value")
	ErrCipherNotFound             = errors.New("cipher not found")
	ErrInvalidCharactersInVideoId = errors.New("invalid characters in video id")
	ErrVideoIdMinLength           = errors.New("the video id must be at least 10 characters long")
	ErrReadOnClosedResBody        = errors.New("http: read on closed response body")
)
