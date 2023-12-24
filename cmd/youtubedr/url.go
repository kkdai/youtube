package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// urlCmd represents the url command
var urlCmd = &cobra.Command{
	Use:   "url",
	Short: "Only output the stream-url to desired video",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		video, format, err := getVideoWithFormat(args[0])
		exitOnError(err)

		url, err := downloader.GetStreamURL(video, format)
		exitOnError(err)

		fmt.Println(url)
	},
}

func init() {
	addVideoSelectionFlags(urlCmd.Flags())
	rootCmd.AddCommand(urlCmd)
}
