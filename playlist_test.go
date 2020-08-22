package youtube

import (
	"testing"
)

func TestYoutube_extractPlaylistID(t *testing.T) {
	attempts := []struct {
		name      string
		url       string
		correctID string
		wantErr   bool
	}{
		{
			"pass-1",
			"https://www.youtube.com/watch?v=9UL390els7M&list=PLqAfPOrmacr963ATEroh67fbvjmTzTEx5",
			"PLqAfPOrmacr963ATEroh67fbvjmTzTEx5",
			false,
		},
		{
			"pass-2",
			"PLqAfPOrmacr963ATEroh67fbvjmTzTEx5", "PLqAfPOrmacr963ATEroh67fbvjmTzTEx5",
			false,
		},
		{
			"pass-3",
			"&list=PLqAfPOrmacr963ATEroh67fbvjmTzTEx5", "PLqAfPOrmacr963ATEroh67fbvjmTzTEx5",
			false,
		},
		{
			"fail-1",
			"https://www.youtube.com/watch?v=9UL390els7M&list=PLqAfPOrmacr963ATEroh67fbvjmTzTEx", "",
			true,
		},
		{
			"fail-2",
			"", "",
			true,
		},
		{
			"fail-3",
			"awevqevqwev", "",
			true,
		},
	}

	for _, v := range attempts {
		ans, err := extractPlaylistID(v.url)

		if err != nil && !v.wantErr {
			t.Errorf("test: %s\nerror: %v\n", v.name, err)
			return
		}

		if err == nil && v.wantErr {
			t.Errorf("Test %s should have errored, url: %v\n", v.name, v.url)
			return
		}

		if ans != v.correctID {
			t.Errorf("Test %s wanted id: %v, entered url: %v\n", v.name, v.correctID, v.url)
		}
	}
}
