package youtube

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testClient = Client{Debug: true}
)

const dwlURL string = "https://www.youtube.com/watch?v=rFejpH_tAHM"
const errURL string = "https://www.youtube.com/watch?v=I8oGsuQ"

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
			expectedErr: ErrInvalidCharactersInVideoId,
		},
		{
			name: "video id is less than 10 characters",
			args: args{
				"rFejpH",
			},
			wantErr:     true,
			expectedErr: ErrVideoIdMinLength,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := extractVideoID(tt.args.url); (err != nil) != tt.wantErr || err != tt.expectedErr {
				t.Errorf("extractVideoID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
