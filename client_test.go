package youtube

import (
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

var testClient = Client{Debug: true}

const (
	dwlURL    string = "https://www.youtube.com/watch?v=rFejpH_tAHM"
	streamURL string = "https://www.youtube.com/watch?v=5qap5aO4i9A"
	errURL    string = "https://www.youtube.com/watch?v=I8oGsuQ"
)

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
	require.NoError(err)
	require.NotNil(video)

	assert.NotEmpty(video.Thumbnails)
	assert.Greater(len(video.Thumbnails), 0)
	assert.NotEmpty(video.Thumbnails[0].URL)
	assert.Empty(video.HLSManifestURL)
	assert.Empty(video.DASHManifestURL)

	assert.Equal("rFejpH_tAHM", video.ID)
	assert.Equal("dotGo 2015 - Rob Pike - Simplicity is Complicated", video.Title)
	assert.Equal("dotconferences", video.Author)
	assert.Equal(1392*time.Second, video.Duration)
	assert.Contains(video.Description, "Go is often described as a simple language.")
	assert.Equal("2015-12-02 00:00:00 +0000 UTC", video.PublishDate.String())
}

func TestGetVideoWithManifestURL(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	video, err := testClient.GetVideo(streamURL)
	require.NoError(err)
	require.NotNil(video)

	assert.NotEmpty(video.Thumbnails)
	assert.Greater(len(video.Thumbnails), 0)
	assert.NotEmpty(video.Thumbnails[0].URL)
	assert.NotEmpty(video.HLSManifestURL)
	assert.NotEmpty(video.DASHManifestURL)
}

func TestGetStream(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	video, err := testClient.GetVideo("https://www.youtube.com/watch?v=BaW_jenozKc")
	require.NoError(err)
	require.NotNil(video)
	require.Greater(len(video.Formats), 0)

	reader, size, err := testClient.GetStream(video, &video.Formats[0])
	require.NoError(err)
	assert.EqualValues(2208750, size)

	data, err := io.ReadAll(reader)
	require.NoError(err)
	assert.Len(data, int(size))
}

func TestGetPlaylist(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	playlist, err := testClient.GetPlaylist("https://www.youtube.com/playlist?list=PL2F4AF82A41D0D2C6")
	require.NoError(err)
	require.NotNil(playlist)

	if assert.Len(playlist.Videos, 1) {
		assert.Equal("lDc63xu98cE", playlist.Videos[0].ID)
		assert.Equal("Dazian Showroom 2010", playlist.Videos[0].Title)
		assert.Equal("Dazian Creative Fabric", playlist.Videos[0].Author)
		assert.Equal(173*time.Second, playlist.Videos[0].Duration)
	}
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
