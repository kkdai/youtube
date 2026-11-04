package youtube

import (
	"errors"
	"fmt"
	"net/url"
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

func TestErrRateLimited(t *testing.T) {
	t.Parallel()

	assert.ErrorIs(t, ErrUnexpectedStatusCode(429), ErrRateLimited)
	assert.ErrorIs(t, fmt.Errorf("wrapped: %w", ErrUnexpectedStatusCode(429)), ErrRateLimited)
	assert.NotErrorIs(t, ErrUnexpectedStatusCode(404), ErrRateLimited)
}

func TestIsRateLimitRedirect(t *testing.T) {
	t.Parallel()

	tooManyRedirects := errors.New("stopped after 10 redirects")

	assert.True(t, isRateLimitRedirect(&url.Error{
		Op:  "Get",
		URL: "https://www.google.com/sorry/index?continue=https://www.youtube.com/",
		Err: tooManyRedirects,
	}))
	assert.False(t, isRateLimitRedirect(&url.Error{
		Op:  "Get",
		URL: "https://www.youtube.com/watch?v=ezfPHVOf-NE",
		Err: tooManyRedirects,
	}))
	assert.False(t, isRateLimitRedirect(nil))
}
