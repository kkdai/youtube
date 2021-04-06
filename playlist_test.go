package youtube

import (
	"testing"
)

func TestYoutube_extractPlaylistID(t *testing.T) {
	attempts := []struct {
		name      string
		url       string
		correctID string
		expected  error
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

	for _, v := range attempts {
		ans, err := extractPlaylistID(v.url)

		if err != v.expected {
			t.Errorf("test: %s\nerror: %v\nexpected: %v\n", v.name, err, v.expected)
			return
		}

		if ans != v.correctID {
			t.Errorf("Test %s wanted id: %v, entered url: %v\n", v.name, v.correctID, v.url)
		}
	}
}
