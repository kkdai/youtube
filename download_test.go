// +build integration

package youtube

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadFromYT_AssignOutputFileName(t *testing.T) {
	y := NewYoutube(false)
	if y == nil {
		t.Error("Cannot init object.")
		return
	}

	// url from issue #25
	testVideoUrl := "https://www.youtube.com/watch?v=54e6lBE3BoQ"
	if err := y.DecodeURL(testVideoUrl); err != nil {
		t.Error("Cannot decode download url")
		return
	}
	curDir, _ := os.Getwd()
	outputDir := curDir + "\\" + downloadToDir
	testcases := []struct {
		name       string
		outputFile string
		itagNo     int
		quality    string
	}{
		{
			name:       "default",
			outputFile: "default_test.mp4",
			itagNo:     0,
			quality:    "",
		},
		{
			name:       "quality:medium",
			outputFile: "medium_test.mp4",
			itagNo:     0,
			quality:    "medium",
		},
	}
	for _, tc := range testcases {
		fmt.Println("download to " + outputDir + "\\" + tc.outputFile)
		t.Run(tc.name, func(t *testing.T) {
			if err := y.StartDownload(outputDir, tc.outputFile, tc.quality, tc.itagNo); err != nil {
				t.Error("Failed in downloading")
			}
		})
	}
}
func TestDownloadFromYT_NoOutputFileName(t *testing.T) {
	y := NewYoutube(false)
	if y == nil {
		t.Error("Cannot init object.")
		return
	}

	// url from issue #21
	testVideoUrl := "https://www.youtube.com/watch?v=n3kPvBCYT3E"
	if err := y.DecodeURL(testVideoUrl); err != nil {
		t.Error("Cannot decode download url")
		return
	}
	curDir, _ := os.Getwd()
	outputDir := filepath.Join(curDir, downloadToDir)
	fmt.Println("download to " + outputDir + "\\" + "Silhouette Eurobeat Remix")
	if err := y.StartDownload(outputDir, "", "", 0); err != nil {
		t.Error("Failed in downloading")
		return
	}
}

func TestDownloadFromYT_HighQuality(t *testing.T) {
	y := NewYoutube(false)
	if y == nil {
		t.Error("Cannot init object.")
		return
	}

	// url from issue #21
	testVideoId := "n3kPvBCYT3E"
	if err := y.DecodeURL(testVideoId); err != nil {
		t.Error("Cannot decode download url")
		return
	}
	curDir, _ := os.Getwd()
	outputDir := filepath.Join(curDir, downloadToDir)
	fmt.Println("download to " + outputDir + "\\" + "Silhouette Eurobeat Remix")
	if err := y.StartDownloadWithHighQuality(outputDir, "", "hd1080"); err != nil {
		t.Error("Failed in downloading")
		return
	}
}

func TestDownloadFromYT_WithItag(t *testing.T) {
	y := NewYoutube(false)
	if y == nil {
		t.Error("Cannot init object.")
		return
	}

	// url from issue #25
	testVideoUrl := "https://www.youtube.com/watch?v=54e6lBE3BoQ"
	if err := y.DecodeURL(testVideoUrl); err != nil {
		t.Error("Cannot decode download url")
		return
	}

	curDir, _ := os.Getwd()
	outputDir := curDir + "\\" + downloadToDir
	testcases := []struct {
		name       string
		outputFile string
		itagNo     int
	}{
		{
			name:       "Format",
			outputFile: "muxedstream_test.mp4",
			itagNo:     18,
		},
		{
			name:       "AdaptiveFormat_video",
			outputFile: "adaptiveStream_video_test.m4v",
			itagNo:     134,
		},
		{
			name:       "AdaptiveFormat_audio",
			outputFile: "adaptiveStream_audio_test.m4a",
			itagNo:     140,
		},
	}

	for _, ts := range testcases {
		t.Run(ts.name, func(t *testing.T) {
			if err := y.StartDownload(outputDir, ts.outputFile, "", ts.itagNo); err != nil {
				t.Errorf("Failed in downloading, err:%v\n", err)
				return
			}
		})
	}
}
