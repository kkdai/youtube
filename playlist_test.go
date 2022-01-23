package youtube

import (
	"encoding/json"
	"os"
	"testing"
	"time"

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
			"pass-5",
			"https://www.youtube.com/watch?v=-T4THwne8IE&list=RD-T4THwne8IE",
			"RD-T4THwne8IE",
			nil,
		},
		{
			"fail-1-playlist-id-44-char",
			"https://www.youtube.com/watch?v=9UL390els7M&list=PLqAfPOrmacr963ATEroh67fbvjmTzTEx5X1212404244", "",
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
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			id, err := extractPlaylistID(v.url)

			assert.Equal(t, v.expectedError, err)
			assert.Equal(t, v.expectedID, id)
		})
	}
}

func TestExtractPlaylist(t *testing.T) {
	f, err := os.Open(testPlaylistResponseDataFile)
	assert.NoError(t, err)
	defer f.Close()
	data, err := extractPlaylistJSON(f)
	assert.NoError(t, err)

	p := &Playlist{ID: testPlaylistID}
	err = json.Unmarshal(data, p)
	assert.NoError(t, err)
	assert.Equal(t, p.Title, "Test Playlist")
	assert.Equal(t, p.Description, "")
	assert.Equal(t, p.Author, "GoogleVoice")
	assert.Equal(t, len(p.Videos), 8)

	v := p.Videos[7]
	assert.Equal(t, v.ID, "dsUXAEzaC3Q")
	assert.Equal(t, v.Title, "Michael Jackson - Bad (Shortened Version)")
	assert.Equal(t, v.Author, "Michael Jackson")
	assert.Equal(t, v.Duration, 4*time.Minute+20*time.Second)

	assert.NotEmpty(t, v.Thumbnails)
	assert.NotEmpty(t, v.Thumbnails[0].URL)
}
