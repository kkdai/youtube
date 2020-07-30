package downloader

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kkdai/youtube"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
)

// Downloader offers high level functions to download videos into files
type Downloader struct {
	Client    youtube.Client
	OutputDir string // optional directory to store the files
}

func (dl *Downloader) getOutputFile(v *youtube.Video, stream *youtube.Stream, outputFile string) (string, error) {

	if outputFile == "" {
		outputFile = SanitizeFilename(v.Title)
		outputFile += pickIdealFileExtension(stream.MimeType)
	}

	if dl.OutputDir != "" {
		if err := os.MkdirAll(dl.OutputDir, 0755); err != nil {
			return "", err
		}
		outputFile = filepath.Join(dl.OutputDir, outputFile)
	}

	return outputFile, nil
}

//Download : Starting download video by arguments.
func (dl *Downloader) Download(ctx context.Context, outputFile string, v *youtube.Video, stream *youtube.Stream) error {
	destFile, err := dl.getOutputFile(v, stream, outputFile)
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
	return dl.videoDLWorker(ctx, out, v, stream)
}

//DownloadWithHighQuality : Starting downloading video with high quality (>720p).
func (dl *Downloader) DownloadWithHighQuality(ctx context.Context, outputFile string, v *youtube.Video, quality string) error {
	var videoStream, audioStream *youtube.Stream

	switch quality {
	case "hd1080":
		videoStream = v.FindStreamByItag(137)
		audioStream = v.FindStreamByItag(140)
	default:
		return fmt.Errorf("unknown quality: %s", quality)
	}

	if videoStream == nil {
		return fmt.Errorf("no Stream video/mp4 for %s found", quality)
	}
	if audioStream == nil {
		return fmt.Errorf("no Stream audio/mp4 for %s found", quality)
	}

	destFile, err := dl.getOutputFile(v, videoStream, outputFile)
	if err != nil {
		return err
	}
	outputDir := filepath.Dir(destFile)

	// Create temporary video file
	videoFile, err := ioutil.TempFile(outputDir, "youtube_*.m4v")
	if err != nil {
		return err
	}
	defer os.Remove(videoFile.Name())

	// Create temporary audio file
	audioFile, err := ioutil.TempFile(outputDir, "youtube_*.m4a")
	if err != nil {
		return err
	}
	defer os.Remove(audioFile.Name())

	dl.logf("Downloading video file...")
	err = dl.videoDLWorker(ctx, videoFile, v, videoStream)
	if err != nil {
		return err
	}

	dl.logf("Downloading audio file...")
	err = dl.videoDLWorker(ctx, audioFile, v, audioStream)
	if err != nil {
		return err
	}

	ffmpegVersionCmd := exec.Command("ffmpeg", "-y",
		"-i", videoFile.Name(),
		"-i", audioFile.Name(),
		"-strict",
		"-2",
		"-shortest",
		destFile,
		"-loglevel", "warning",
	)
	ffmpegVersionCmd.Stderr = os.Stderr
	ffmpegVersionCmd.Stdout = os.Stdout
	dl.logf("merging video and audio to %s", destFile)

	return ffmpegVersionCmd.Run()
}

func (dl *Downloader) videoDLWorker(ctx context.Context, out *os.File, video *youtube.Video, stream *youtube.Stream) error {
	resp, err := dl.Client.Download(ctx, video, stream)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	prog := &progress{
		contentLength: float64(resp.ContentLength),
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
	reader := bar.ProxyReader(resp.Body)
	defer dl.close(reader)

	mw := io.MultiWriter(out, prog)
	_, err = io.Copy(mw, reader)
	if err != nil {
		return err
	}
	progress.Wait()
	return nil
}

func (dl *Downloader) close(r io.ReadCloser) {
	_, err := io.Copy(ioutil.Discard, r)
	if err != nil && err.Error() != youtube.ErrReadOnClosedResBody.Error() {
		dl.logf("failed to exhaust reader: %s", err)
	}
	err = r.Close()
	if err != nil {
		dl.logf("response close err %s", err)
	}
}

func (dl *Downloader) logf(format string, v ...interface{}) {
	if dl.Client.Debug {
		log.Printf(format, v...)
	}
}
