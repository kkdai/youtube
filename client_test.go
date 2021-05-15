package youtube

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	assert.EqualError(t, err, `response status: 'fail', reason: 'Invalid parameters.'`)
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
