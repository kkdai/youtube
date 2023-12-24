package youtube

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		err      error
		expected string
	}{
		{ErrUnexpectedStatusCode(404), "unexpected status code: 404"},
		{ErrPlayabiltyStatus{"invalid", "for that reason"}, "cannot playback and download, status: invalid, reason: for that reason"},
		{ErrPlaylistStatus{"for that reason"}, "could not load playlist: for that reason"},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.EqualError(t, tt.err, tt.expected)
		})
	}
}
