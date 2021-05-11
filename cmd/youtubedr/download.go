package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:     "download",
	Short:   "Downloads a video from youtube",
	Example: `youtubedr -o "Campaign Diary".mp4 https://www.youtube.com/watch\?v\=XbNghLqsVwU`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		exitOnError(download(args[0]))
	},
}

var (
	ffmpegCheck error
	outputFile  string
	outputDir   string
)

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVarP(&outputFile, "filename", "o", "", "The output file, the default is genated by the video title.")
	downloadCmd.Flags().StringVarP(&outputDir, "directory", "d", ".", "The output directory.")
	addQualityFlag(downloadCmd.Flags())
	addMimeTypeFlag(downloadCmd.Flags())
}

func download(id string) error {
	video, format, err := getVideoWithFormat(id)
	if err != nil {
		downloader.Logf("⛔ %s: '%s'\n", id, err)
		return err
	}
	downloader.Logf("▶ %s: '%s'\n", video.ID, video.Title)
	downloader.Logf("download to directory: %s\n", outputDir)

	if strings.HasPrefix(outputQuality, "hd") {
		if err := checkFFMPEG(); err != nil {
			return err
		}
		return downloader.DownloadComposite(context.Background(), outputFile, video, outputQuality, mimetype)
	}

	return downloader.Download(context.Background(), video, format, outputFile)
}

func checkFFMPEG() error {
	downloader.Logf("check ffmpeg is installed....")
	if err := exec.Command("ffmpeg", "-version").Run(); err != nil {
		ffmpegCheck = fmt.Errorf("please check ffmpegCheck is installed correctly")
	}

	return ffmpegCheck
}
