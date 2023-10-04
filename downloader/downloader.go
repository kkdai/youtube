package downloader

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/takefy-dev/youtube"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
)

// Downloader offers high level functions to download videos into files
type Downloader struct {
	youtube.Client
	OutputDir string // optional directory to store the files
}

func (dl *Downloader) getOutputFile(v *youtube.Video, format *youtube.Format, outputFile string) (string, error) {
	if outputFile == "" {
		outputFile = SanitizeFilename(v.Title)
		outputFile += pickIdealFileExtension(format.MimeType)
	}

	if dl.OutputDir != "" {
		if err := os.MkdirAll(dl.OutputDir, 0o755); err != nil {
			return "", err
		}
		outputFile = filepath.Join(dl.OutputDir, outputFile)
	}

	return outputFile, nil
}

// Download : Starting download video by arguments.
func (dl *Downloader) Download(ctx context.Context, v *youtube.Video, format *youtube.Format, outputFile string) error {
	dl.logf("Video '%s' - Quality '%s' - Codec '%s'", v.Title, format.QualityLabel, format.MimeType)
	destFile, err := dl.getOutputFile(v, format, outputFile)
	if err != nil {
		return err
	}

	// Create output file
	out, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer out.Close()

	dl.logf("Download to file=%s", destFile)
	return dl.videoDLWorker(ctx, out, v, format)
}

// DownloadComposite : Downloads audio and video streams separately and merges them via ffmpeg.
func (dl *Downloader) DownloadComposite(ctx context.Context, outputFile string, v *youtube.Video, quality string, mimetype string) error {
	videoFormat, audioFormat, err1 := getVideoAudioFormats(v, quality, mimetype)
	if err1 != nil {
		return err1
	}

	dl.logf("Video '%s' - Quality '%s' - Video Codec '%s' - Audio Codec '%s'", v.Title, videoFormat.QualityLabel, videoFormat.MimeType, audioFormat.MimeType)
	destFile, err := dl.getOutputFile(v, videoFormat, outputFile)
	if err != nil {
		return err
	}
	outputDir := filepath.Dir(destFile)

	// Create temporary video file
	videoFile, err := os.CreateTemp(outputDir, "youtube_*.m4v")
	if err != nil {
		return err
	}
	defer os.Remove(videoFile.Name())

	// Create temporary audio file
	audioFile, err := os.CreateTemp(outputDir, "youtube_*.m4a")
	if err != nil {
		return err
	}
	defer os.Remove(audioFile.Name())

	dl.logf("Downloading video file...")
	err = dl.videoDLWorker(ctx, videoFile, v, videoFormat)
	if err != nil {
		return err
	}

	dl.logf("Downloading audio file...")
	err = dl.videoDLWorker(ctx, audioFile, v, audioFormat)
	if err != nil {
		return err
	}

	//nolint:gosec
	ffmpegVersionCmd := exec.Command("ffmpeg", "-y",
		"-i", videoFile.Name(),
		"-i", audioFile.Name(),
		"-c", "copy", // Just copy without re-encoding
		"-shortest", // Finish encoding when the shortest input stream ends
		destFile,
		"-loglevel", "warning",
	)
	ffmpegVersionCmd.Stderr = os.Stderr
	ffmpegVersionCmd.Stdout = os.Stdout
	dl.logf("merging video and audio to %s", destFile)

	return ffmpegVersionCmd.Run()
}

func getVideoAudioFormats(v *youtube.Video, quality string, mimetype string) (*youtube.Format, *youtube.Format, error) {
	var videoFormat, audioFormat *youtube.Format
	var videoFormats, audioFormats youtube.FormatList

	formats := v.Formats
	if mimetype != "" {
		formats = formats.Type(mimetype)
	}

	videoFormats = formats.Type("video").AudioChannels(0)
	audioFormats = formats.Type("audio")

	if quality != "" {
		videoFormats = videoFormats.Quality(quality)
	}

	if len(videoFormats) > 0 {
		videoFormats.Sort()
		videoFormat = &videoFormats[0]
	}

	if len(audioFormats) > 0 {
		audioFormats.Sort()
		for _, format := range audioFormats {
			if format.AudioTrack.Id == v.DefaultAudioTrackId {
				audioFormat = &format
				break
			}
		}
	}

	if videoFormat == nil {
		return nil, nil, fmt.Errorf("no video format found after filtering")
	}

	if audioFormat == nil {
		return nil, nil, fmt.Errorf("no audio format found after filtering")
	}

	return videoFormat, audioFormat, nil
}

func (dl *Downloader) videoDLWorker(ctx context.Context, out *os.File, video *youtube.Video, format *youtube.Format) error {
	stream, size, err := dl.GetStreamContext(ctx, video, format)
	if err != nil {
		return err
	}

	prog := &progress{
		contentLength: float64(size),
	}

	// create progress bar
	progress := mpb.New(mpb.WithWidth(64))
	bar := progress.AddBar(
		int64(prog.contentLength),

		mpb.PrependDecorators(
			decor.CountersKibiByte("% .2f / % .2f"),
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_GO, 90),
			decor.Name(" ] "),
			decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
		),
	)

	reader := bar.ProxyReader(stream)
	mw := io.MultiWriter(out, prog)
	_, err = io.Copy(mw, reader)
	if err != nil {
		return err
	}

	progress.Wait()
	return nil
}

func (dl *Downloader) logf(format string, v ...interface{}) {
	if dl.Debug {
		log.Printf(format, v...)
	}
}
