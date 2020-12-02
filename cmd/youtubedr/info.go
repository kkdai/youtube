package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print metadata of the desired video",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		video, err := getDownloader().GetVideo(args[0])
		exitOnError(err)

		fmt.Println("Title:      ", video.Title)
		fmt.Println("Author:     ", video.Author)
		fmt.Println("Duration:   ", video.Duration)
		fmt.Println("Description:", video.Description)
		fmt.Println()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoWrapText(false)
		table.SetHeader([]string{"itag", "video quality", "audio quality", "size [MB]", "bitrate", "MimeType"})

		for _, format := range video.Formats {
			bitrate := format.AverageBitrate
			if bitrate == 0 {
				// Some formats don't have the average bitrate
				bitrate = format.Bitrate
			}

			size, _ := strconv.ParseInt(format.ContentLength, 10, 64)
			if size == 0 {
				// Some formats don't have this information
				size = int64(float64(bitrate) * video.Duration.Seconds() / 8)
			}

			table.Append([]string{
				strconv.Itoa(format.ItagNo),
				format.QualityLabel,
				strings.ToLower(strings.TrimPrefix(format.AudioQuality, "AUDIO_QUALITY_")),
				fmt.Sprintf("%0.1f", float64(size)/1024/1024),
				strconv.Itoa(bitrate),
				format.MimeType,
			})
		}
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
