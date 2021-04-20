package main

import (
	"fmt"
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type PlaylistInfo struct {
	Title  string
	Author string
	Videos []VideoInfo
}

var (
	playlistOutputWriters = map[string]outputWriter{
		"json": jsonOutput(),
		"xml":  xmlOutput(),
		"plain": func(info interface{}, w io.Writer) error {
			i, ok := info.(PlaylistInfo)
			if !ok {
				return fmt.Errorf("input is not PlaylistInfo")
			}
			fmt.Println("Title:      ", i.Title)
			fmt.Println("Author:     ", i.Author)
			fmt.Println("# Videos:   ", len(i.Videos))

			fmt.Println()

			table := tablewriter.NewWriter(os.Stdout)
			table.SetAutoWrapText(false)
			table.SetHeader([]string{"ID", "Author", "Title", "Duration"})

			for _, vid := range i.Videos {
				table.Append([]string{
					vid.ID,
					vid.Author,
					vid.Title,
					fmt.Sprintf("%v", vid.Duration),
				})
			}
			table.Render()
			return nil
		},
	}
	// listCmd represents the list command
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "Print metadata of the desired playlist",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			outputFunc = playlistOutputWriters[outputFormat]
			if outputFunc == nil {
				return fmt.Errorf("output format %s is not valid", outputFormat)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			playlist, err := getDownloader().GetPlaylist(args[0])
			exitOnError(err)

			playlistInfo := PlaylistInfo{
				Title:  playlist.Title,
				Author: playlist.Author,
			}
			for _, v := range playlist.Videos {
				playlistInfo.Videos = append(playlistInfo.Videos, VideoInfo{
					ID:       v.ID,
					Title:    v.Title,
					Author:   v.Author,
					Duration: v.Duration.String(),
				})
			}
			exitOnError(outputFunc(playlistInfo, os.Stdout))
		},
	}
)

func init() {
	rootCmd.AddCommand(listCmd)
	addFormatFlag(listCmd.Flags())
}
