package youtube

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYoutube_extractPlaylistID(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		expectedID    string
		expectedError error
	}{
		{
			"pass-1",
			"https://www.youtube.com/watch?v=9UL390els7M&list=PLqAfPOrmacr963ATEroh67fbvjmTzTEx5",
			"PLqAfPOrmacr963ATEroh67fbvjmTzTEx5",
			nil,
		},

		{
			"pass-2",
			"PLqAfPOrmacr963ATEroh67fbvjmTzTEx5",
			"PLqAfPOrmacr963ATEroh67fbvjmTzTEx5",
			nil,
		},
		{
			"pass-3",
			"&list=PLqAfPOrmacr963ATEroh67fbvjmTzTEx5",
			"PLqAfPOrmacr963ATEroh67fbvjmTzTEx5",
			nil,
		},
		{
			"pass-4 (extra params)",
			"https://www.youtube.com/watch?v=9UL390els7M&list=PLqAfPOrmacr963ATEroh67fbvjmTzTEx5&foo=bar&baz=babar",
			"PLqAfPOrmacr963ATEroh67fbvjmTzTEx5",
			nil,
		},
		{
			"fail-1",
			"https://www.youtube.com/watch?v=9UL390els7M&list=PLqAfPOrmacr963ATEroh67fbvjmTzTEx5X", "",
			ErrInvalidPlaylist,
		},
		{
			"fail-2",
			"", "",
			ErrInvalidPlaylist,
		},
		{
			"fail-3",
			"awevqevqwev", "",
			ErrInvalidPlaylist,
		},
		{
			"fail-4",
			"https://www.youtube.com/watch?v=9UL390els7M&list=PLqAfPOrmacr963A&foo=bar", "",
			ErrInvalidPlaylist,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			id, err := extractPlaylistID(v.url)

			assert.Equal(t, v.expectedError, err)
			assert.Equal(t, v.expectedID, id)
		})
	}
}
