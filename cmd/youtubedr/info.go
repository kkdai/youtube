package main

import (
	"fmt"
	"os"
	"strconv"

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
		table.SetHeader([]string{"itag", "quality", "MimeType"})

		for _, itag := range video.Formats {
			table.Append([]string{strconv.Itoa(itag.ItagNo), itag.Quality, itag.MimeType})
		}
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
