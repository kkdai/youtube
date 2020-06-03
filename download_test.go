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
	outputFile := "download_test.mp4"
	fmt.Println("download to " + outputDir + "\\" + outputFile)
	if err := y.StartDownload(outputDir, outputFile, "", 0); err != nil {
		t.Error("Failed in downloading")
		return
	}
}
func TestDownloadFromYT_NoOutputFileName(t *testing.T) {
	y := NewYoutube(false)
	if y == nil {
		t.Error("Cannot init object.")
		return
	}

	// url from issue #
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
	outputFile := "download_test.mp4"
	if err := y.StartDownload(outputDir, outputFile, "", 18); err != nil {
		t.Error("Failed in downloading")
		return
	}
}
