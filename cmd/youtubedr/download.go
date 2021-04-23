package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/kkdai/youtube/v2"
	"github.com/spf13/cobra"
)

var (
	errNotID = fmt.Errorf("cannot detect ID in given input")
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:     "download",
	Short:   "Downloads a video or a playlist from youtube",
	Example: `youtubedr -o "Campaign Diary".mp4 https://www.youtube.com/watch\?v\=XbNghLqsVwU`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if playlistID, err := youtube.ExtractPlaylistID(args[0]); err != nil {
			videoID, err := youtube.ExtractVideoID(args[0])
			if err != nil {
				exitOnError(errNotID)
			}
			log.Printf(
				"download video %s to directory %s\n",
				videoID,
				outputDir,
			)
			exitOnError(download(videoID))
		} else {
			playlist, err := getDownloader().GetPlaylist(playlistID)
			if err != nil {
				exitOnError(err)
			}
			log.Printf(
				"download %d videos from playlist %s to directory %s\n",
				len(playlist.Videos),
				playlist.ID,
				outputDir,
			)
			var errors []error
			outputFileOrigin := outputFile
			for i, v := range playlist.Videos {
				if len(outputFileOrigin) != 0 {
					outputFile = fmt.Sprintf("%s-%d", outputFile, i)
				}
				if err := download(v.ID); err != nil {
					errors = append(errors, err)
				}
			}
			if len(errors) > 0 {
				exitOnErrors(errors)
			}
		}
	},
}

var (
	ffmpegCheck error
	outputFile  string
	outputDir   string
)

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVarP(&outputFile, "filename", "o", "", "The output file, the default is generated from the video title.")
	downloadCmd.Flags().StringVarP(&outputDir, "directory", "d", ".", "The output directory.")
	addQualityFlag(downloadCmd.Flags())
	addMimeTypeFlag(downloadCmd.Flags())
}

func download(id string) error {
	video, format, err := getVideoWithFormat(id)
	if err != nil {
		return err
	}

	if strings.HasPrefix(outputQuality, "hd") {
		if err := checkFFMPEG(); err != nil {
			return err
		}
		return downloader.DownloadComposite(context.Background(), outputFile, video, outputQuality, mimetype)
	}

	return downloader.Download(context.Background(), video, format, outputFile)
}

func checkFFMPEG() error {
	fmt.Println("check ffmpeg is installed....")
	if err := exec.Command("ffmpeg", "-version").Run(); err != nil {
		ffmpegCheck = fmt.Errorf("please check ffmpegCheck is installed correctly")
	}

	return ffmpegCheck
}
