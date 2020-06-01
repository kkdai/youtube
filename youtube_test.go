package youtube

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
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
	}{
		{name: "Default"},
		{name: "with outputDir", outputDir: dfPath},
		{name: "SpecificQuality", quality: "hd720"},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			y := NewYoutube(false)
			if y == nil {
				t.Error("Cannot init object.")
				return
			}

			if err := y.StartDownload(tc.outputDir, tc.ouputFile, tc.quality); err == nil {
				t.Error("No video URL input should not download.")
				return
			}
		})
	}
}

func TestDownloadSpecificItag(t *testing.T) {
	y := NewYoutube(false)
	if y == nil {
		t.Error("Cannot init object.")
		return
	}

	if err := y.StartDownloadWithItag(dfPath, 22); err == nil {
		t.Error("No video URL input should not download.")
		return
	}
}

func TestDownloadErrEmptyStreamList(t *testing.T) {
	y := NewYoutube(false)
	if y == nil {
		t.Error("Cannot init object.")
		return
	}

	if err := y.StartDownload("", "", ""); err != ErrEmptyStreamList {
		t.Error("no error returned by empty stream list")
		return
	}
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
