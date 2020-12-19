package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// urlCmd represents the url command
var urlCmd = &cobra.Command{
	Use:   "url",
	Short: "Only output the stream-url to desired video",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var errors []string
		for _, videoURL := range args {
			video, format, err := getVideoWithFormat(videoURL)
			if err != nil {
				errors = append(errors, err.Error())
				continue
			}

			fmt.Printf("Video '%s' - Quality '%s' - Codec '%s'", video.Title, format.QualityLabel, format.MimeType)
			url, err := downloader.GetStreamURL(video, format)
			if err != nil {
				errors = append(errors, err.Error())
			}

			fmt.Println(url)
		}
		if len(errors) > 0 {
			return fmt.Errorf(strings.Join(errors, " | "))
		}
		return nil
	},
}

func init() {
	addQualityFlag(urlCmd.Flags())
	addCodecFlag(urlCmd.Flags())
	rootCmd.AddCommand(urlCmd)
}
