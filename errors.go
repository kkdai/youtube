package youtube

import "errors"

var (
	ErrEmptyStreamList = errors.New("Empty stream list")
	ErrItagNotFound    = errors.New("Invalid itag value, please specify correct value.")
)
