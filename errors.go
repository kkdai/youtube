package youtube

import (
	"errors"
	"fmt"
)

var (
	ErrCipherNotFound             = errors.New("cipher not found")
	ErrInvalidCharactersInVideoID = errors.New("invalid characters in video id")
	ErrVideoIDMinLength           = errors.New("the video id must be at least 10 characters long")
	ErrReadOnClosedResBody        = errors.New("http: read on closed response body")
)

type ErrResponseStatus struct {
	Status string
	Reason string
}

func (err ErrResponseStatus) Error() string {
	if err.Status == "" {
		return "no response status found in the server's answer"
	}

	if err.Reason == "" {
		return fmt.Sprintf("response status: '%s', no reason given", err.Status)
	}

	return fmt.Sprintf("response status: '%s', reason: '%s'", err.Status, err.Reason)
}

type ErrPlayabiltyStatus struct {
	Status string
	Reason string
}

func (err ErrPlayabiltyStatus) Error() string {
	return fmt.Sprintf("cannot playback and download, status: %s, reason: %s", err.Status, err.Reason)
}

// ErrUnexpectedStatusCode is returned on unexpected HTTP status codes
type ErrUnexpectedStatusCode int

func (err ErrUnexpectedStatusCode) Error() string {
	return fmt.Sprintf("unexpected status code: %d", err)
}
