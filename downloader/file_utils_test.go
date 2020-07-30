package downloader

import (
	"testing"
)

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
