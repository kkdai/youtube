package youtube

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyStreamList = errors.New("Empty stream list")
	ErrItagNotFound    = errors.New("Invalid itag value, please specify correct value.")
	ErrCipherNotFound  = errors.New("cipher not found")
)

type ErrDecodingStreamInfo struct {
	streamPos int
}

func (err ErrDecodingStreamInfo) Error() string {
	return fmt.Sprintf("An error occurred while decoding one of the video's stream's information: stream %d.\n", err.streamPos)
}
