package youtube

import (
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

const (
	dwlURL    string = "https://www.youtube.com/watch?v=rFejpH_tAHM"
	streamURL string = "https://www.youtube.com/watch?v=a9LDPn-MO4I"
	errURL    string = "https://www.youtube.com/watch?v=I8oGsuQ"
)

var testClient = Client{}
var testWebClient = Client{client: &WebClient}

func TestParseVideo(t *testing.T) {
	video, err := testClient.GetVideo(dwlURL)
	assert.NoError(t, err)
	assert.NotNil(t, video)

	_, err = testClient.GetVideo(errURL)
	assert.IsType(t, err, &ErrPlayabiltyStatus{})
}

func TestYoutube_findVideoID(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		expectedErr error
	}{
		{
			name: "valid url",
			args: args{
				dwlURL,
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "valid id",
			args: args{
				"rFejpH_tAHM",
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "invalid character in id",
			args: args{
				"<M13",
			},
			wantErr:     true,
			expectedErr: ErrInvalidCharactersInVideoID,
		},
		{
			name: "video id is less than 10 characters",
			args: args{
				"rFejpH",
			},
			wantErr:     true,
			expectedErr: ErrVideoIDMinLength,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ExtractVideoID(tt.args.url); (err != nil) != tt.wantErr || err != tt.expectedErr {
				t.Errorf("extractVideoID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetVideoWithoutManifestURL(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	video, err := testClient.GetVideo(dwlURL)
	require.NoError(err, "get video")
	require.NotNil(video)

	assert.NotEmpty(video.Thumbnails)
	assert.Greater(len(video.Thumbnails), 0)
	assert.NotEmpty(video.Thumbnails[0].URL)
	assert.Empty(video.HLSManifestURL)
	assert.Empty(video.DASHManifestURL)

	assert.NotEmpty(video.CaptionTracks)
	assert.Greater(len(video.CaptionTracks), 0)
	assert.NotEmpty(video.CaptionTracks[0].BaseURL)

	assert.Equal("rFejpH_tAHM", video.ID)
	assert.Equal("dotGo 2015 - Rob Pike - Simplicity is Complicated", video.Title)
	assert.Equal("dotconferences", video.Author)
	assert.GreaterOrEqual(video.Duration, 1390*time.Second)
	assert.Contains(video.Description, "Go is often described as a simple language.")

	// Publishing date doesn't seem to be present in android client
	// assert.Equal("2015-12-02 00:00:00 +0000 UTC", video.PublishDate.String())
}

func TestWebClientGetVideoWithoutManifestURL(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	video, err := testWebClient.GetVideo(dwlURL)
	require.NoError(err, "get video")
	require.NotNil(video)

	assert.NotEmpty(video.Thumbnails)
	assert.Greater(len(video.Thumbnails), 0)
	assert.NotEmpty(video.Thumbnails[0].URL)
	assert.Empty(video.HLSManifestURL)
	assert.Empty(video.DASHManifestURL)

	assert.NotEmpty(video.CaptionTracks)
	assert.Greater(len(video.CaptionTracks), 0)
	assert.NotEmpty(video.CaptionTracks[0].BaseURL)

	assert.Equal("rFejpH_tAHM", video.ID)
	assert.Equal("dotGo 2015 - Rob Pike - Simplicity is Complicated", video.Title)
	assert.Equal("dotconferences", video.Author)
	assert.GreaterOrEqual(video.Duration, 1390*time.Second)
	assert.Contains(video.Description, "Go is often described as a simple language.")

	// Publishing date and channel handle are present in web client
	//assert.Equal("2015-12-02 00:00:00 +0000 UTC", video.PublishDate.String())

	assert.Equal("@dotconferences", video.ChannelHandle)
}

func TestGetVideoWithManifestURL(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	video, err := testClient.GetVideo(streamURL)
	require.NoError(err)
	require.NotNil(video)

	assert.NotEmpty(video.Formats)
	assert.NotEmpty(video.Thumbnails)
	assert.Greater(len(video.Thumbnails), 0)
	assert.NotEmpty(video.Thumbnails[0].URL)

	format := video.Formats[0]
	r, size, err := testClient.GetStream(video, &format)
	if assert.NoError(err) {
		r.Close()
	}
	assert.NotZero(size)
}

func TestGetVideo_MultiLanguage(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	video, err := testClient.GetVideo("https://www.youtube.com/watch?v=pU9sHwNKc2c")
	require.NoError(err)
	require.NotNil(video)

	// collect languages
	var languageNames, lanaguageIDs []string
	for _, format := range video.Formats {
		if format.AudioTrack != nil {
			languageNames = append(languageNames, format.LanguageDisplayName())
			lanaguageIDs = append(lanaguageIDs, format.AudioTrack.ID)
		}
	}

	assert.Contains(languageNames, "English original")
	assert.Contains(languageNames, "Portuguese (Brazil)")
	assert.Contains(lanaguageIDs, "en.4")
	assert.Contains(lanaguageIDs, "pt-BR.3")

	assert.Empty(video.Formats.Language("Does not exist"))
	assert.NotEmpty(video.Formats.Language("English original"))
}

func TestGetStream(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	expectedSize := 988479

	// Create testclient to enforce re-using of routines
	testClient := Client{
		MaxRoutines: 10,
		ChunkSize:   int64(expectedSize) / 11,
	}

	// Download should not last longer than a minute.
	// Otherwise we assume Youtube is throtteling us.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	video, err := testClient.GetVideoContext(ctx, "https://www.youtube.com/watch?v=BaW_jenozKc")
	require.NoError(err)
	require.NotNil(video)
	require.Greater(len(video.Formats), 0)

	reader, size, err := testClient.GetStreamContext(ctx, video, &video.Formats[0])
	require.NoError(err)
	assert.EqualValues(expectedSize, size)

	data, err := io.ReadAll(reader)
	require.NoError(err)
	assert.Len(data, int(size))
}

func TestGetPlaylist(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	playlist, err := testClient.GetPlaylist("https://www.youtube.com/playlist?list=PL59FEE129ADFF2B12")
	require.NoError(err)
	require.NotNil(playlist)

	assert.Equal(playlist.Title, "Test Playlist")
	assert.Equal(playlist.Description, "")
	assert.Equal(playlist.Author, "GoogleVoice")
	assert.Equal(len(playlist.Videos), 8)

	v := playlist.Videos[7]
	assert.Equal(v.ID, "dsUXAEzaC3Q")
	assert.Equal(v.Title, "Michael Jackson - Bad (Shortened Version)")
	assert.Equal(v.Author, "Michael Jackson")
	assert.Equal(v.Duration, 4*time.Minute+20*time.Second)

	assert.NotEmpty(v.Thumbnails)
	assert.NotEmpty(v.Thumbnails[0].URL)
}

func TestGetBigPlaylist(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	playlist, err := testClient.GetPlaylist("https://www.youtube.com/playlist?list=PLTC7VQ12-9raqhLCx1S1E_ic35t94dj28")
	require.NoError(err)
	require.NotNil(playlist)

	assert.NotEmpty(playlist.Title)
	assert.NotEmpty(playlist.Description)
	assert.NotEmpty(playlist.Author)

	assert.Greater(len(playlist.Videos), 300)
	assert.NotEmpty(playlist.Videos[300].ID)

	t.Logf("Playlist Title: %s, Video Count: %d", playlist.Title, len(playlist.Videos))
}

func TestClient_httpGetBodyBytes(t *testing.T) {
	tests := []struct {
		url           string
		errorContains string
	}{
		{"unknown://", "unsupported protocol scheme"},
		{"invalid\nurl", "invalid control character in URL"},
		{"http://unknown-host/", "dial tcp"},
		{"http://example.com/does-not-exist", "unexpected status code: 404"},
		{"http://example.com/", ""},
	}
	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err := testClient.httpGetBodyBytes(ctx, tt.url)

			if tt.errorContains == "" {
				assert.NoError(t, err)
			} else if assert.Error(t, err) {
				assert.Contains(t, err.Error(), tt.errorContains)
			}
		})
	}
}
