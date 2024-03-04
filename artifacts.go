package youtube

import (
	"log/slog"
	"os"
	"path/filepath"
)

// destination for artifacts, used by integration tests
var artifactsFolder = os.Getenv("ARTIFACTS")

func writeArtifact(name string, content []byte) {
	// Ensure folder exists
	err := os.MkdirAll(artifactsFolder, os.ModePerm)
	if err != nil {
		slog.Error("unable to create artifacts folder", "path", artifactsFolder, "error", err)
		return
	}

	path := filepath.Join(artifactsFolder, name)
	err = os.WriteFile(path, content, 0600)

	log := slog.With("path", path)
	if err != nil {
		log.Error("unable to write artifact", "error", err)
	} else {
		log.Debug("artifact created")
	}
}
