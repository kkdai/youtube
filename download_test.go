// +build integration

package youtube

import (
	"fmt"
	"os"
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
	outputFile := "download_test.mp4"
	fmt.Println("download to " + outputDir + "\\" + outputFile)
	if err := y.StartDownload(outputDir, outputFile); err != nil {
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

	path, _ := os.Getwd()
	path += "\\" + downloadToDir + "\\download_test.mp4"
	fmt.Println("download to " + path)
	if err := y.StartDownloadWithItag(path, 18); err != nil {
		t.Error("Failed in downloading")
		return
	}
}
