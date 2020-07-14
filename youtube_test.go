package youtube

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const dwlURL string = "https://www.youtube.com/watch?v=rFejpH_tAHM"
const errURL string = "https://www.youtube.com/watch?v=I8oGsuQ"
const downloadToDir = "download_test"

var dfPath string

func TestMain(m *testing.M) {
	//init download path
	usr, _ := user.Current()
	dfPath = filepath.Join(usr.HomeDir, "Movies", "test")

	path, _ := os.Getwd()
	path += "\\" + downloadToDir
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatal(err.Error())
	}

	exitCode := m.Run()
	// the following code doesn't work under debugger, please delete download files manually
	if err := os.RemoveAll(path); err != nil {
		log.Fatal(err.Error())
	}
	os.Exit(exitCode)
}

func TestDownload(t *testing.T) {
	testcases := []struct {
		name      string
		outputDir string
		ouputFile string
		quality   string
		itag      int
	}{
		{name: "Default"},
		{name: "with outputDir", outputDir: dfPath},
		{name: "SpecificQuality", quality: "hd720"},
		{name: "SpecificITag", itag: 22},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			y := NewYoutube(false)
			if y == nil {
				t.Error("Cannot init object.")
				return
			}

			if err := y.StartDownload(tc.outputDir, tc.ouputFile, tc.quality, tc.itag); err == nil {
				t.Error("No video URL input should not download.")
				return
			}
		})
	}
}

func TestDownloadError(t *testing.T) {
	y := NewYoutube(false)
	if y == nil {
		t.Error("Cannot init object.")
		return
	}
	t.Run("empty Stream list error", func(t *testing.T) {
		if err := y.StartDownload("", "", "", 0); err != ErrEmptyStreamList {
			t.Error("no err returned for empty Stream list")
		}
	})

	t.Run("itag not found error", func(t *testing.T) {
		y.StreamList = append(y.StreamList, Stream{})
		if err := y.StartDownload("", "", "", 18); err != ErrItagNotFound {
			t.Error("no error returned for itag not found")
		}
	})
}

func TestParseVideo(t *testing.T) {
	y := NewYoutube(false)
	if y == nil {
		t.Error("Cannot init object.")
		return
	}

	if err := y.DecodeURL(dwlURL); err != nil {
		t.Error("This video parsing should work well.")
		return
	}

	if err := y.DecodeURL(errURL); err == nil {
		t.Error("This video parsing should not work well.")
		return
	}
}

func TestSanitizeFilename(t *testing.T) {
	fileName := "a<b>c:d\\e\"f/g\\h|i?j*k"
	sanitized := SanitizeFilename(fileName)
	if sanitized != "abcdefghijk" {
		t.Error("Invalid characters must get stripped")
	}

	fileName = "aB Cd"
	sanitized = SanitizeFilename(fileName)
	if sanitized != "aB Cd" {
		t.Error("Casing and whitespaces must be preserved")
	}

	fileName = "~!@#$%^&()[].,"
	sanitized = SanitizeFilename(fileName)
	if sanitized != "~!@#$%^&()[].," {
		t.Error("The common harmless symbols should remain valid")
	}
}

func TestGetItagInfo(t *testing.T) {
	type args struct {
		StreamList []Stream
	}
	videoQuality := "TestQuality"
	videoType := "TestType"
	videoTitle := "TestTitle"
	videoAuthor := "TestAuthor"
	tests := []struct {
		name string
		args args
		want *ItagInfo
	}{
		{
			name: "no itag",
			args: args{
				StreamList: nil,
			},
			want: nil,
		},
		{
			name: "one itag",
			args: args{
				StreamList: []Stream{
					{
						Quality: videoQuality,
						Type:    videoType,
						URL:     "",
						ItagNo:  0,
					},
				},
			},
			want: &ItagInfo{
				Title:  videoTitle,
				Author: videoAuthor,
				Itags: []Itag{{
					ItagNo:  0,
					Quality: videoQuality,
					Type:    videoType,
				}},
			},
		},
		{
			name: "two itags",
			args: args{
				StreamList: []Stream{
					{
						Quality: videoQuality,
						Type:    videoType,
						URL:     "",
						ItagNo:  0,
					},
					{
						Quality: videoQuality,
						Type:    videoType,
						URL:     "",
						ItagNo:  0,
					},
				},
			},
			want: &ItagInfo{
				Title:  videoTitle,
				Author: videoAuthor,
				Itags: []Itag{
					{
						ItagNo:  0,
						Quality: videoQuality,
						Type:    videoType,
					},
					{
						ItagNo:  0,
						Quality: videoQuality,
						Type:    videoType,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			y := &Youtube{
				StreamList: tt.args.StreamList,
				Author:     videoAuthor,
				Title:      videoTitle,
			}
			if got := y.GetItagInfo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetItagInfo() = %v, want %v", got, tt.want)
			}
		})
	}
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
			y := NewYoutube(false)
			if err := y.findVideoID(tt.args.url); (err != nil) != tt.wantErr || err != tt.expectedErr {
				t.Errorf("findVideoID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestYoutube_StartDownloadWithHighQuality(t *testing.T) {
	tests := []struct {
		name    string
		stream  []Stream
		wantErr bool
		message string
	}{
		{
			name:    "video Stream not found",
			stream:  []Stream{},
			wantErr: true,
			message: "no Stream video/mp4",
		},
		{
			name:    "audio Stream not found",
			stream:  []Stream{{ItagNo: 137}},
			wantErr: true,
			message: "no Stream audio/mp4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			y := NewYoutube(false)

			if err := y.StartDownloadWithHighQuality("", "", "hd1080"); (err != nil) != tt.wantErr && !strings.Contains(err.Error(), tt.message) {
				t.Errorf("StartDownloadWithHighQuality() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestYoutube_getStreamUrl(t *testing.T) {
	type args struct {
		stream Stream
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name: "url is not empty",
			args: args{
				stream: Stream{
					URL: "test",
				},
			},
			want:    "test",
			wantErr: nil,
		},
		{
			name: "url and cipher is empty",
			args: args{
				stream: Stream{
					URL:    "",
					Cipher: "",
				},
			},
			want:    "",
			wantErr: ErrCipherNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			y := NewYoutube(false)
			got, err := y.getStreamUrl(tt.args.stream)
			if err != tt.wantErr {
				t.Errorf("getStreamUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getStreamUrl() got = %v, want %v", got, tt.want)
			}
		})
	}
}
