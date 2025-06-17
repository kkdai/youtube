package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type PlaylistInfo struct {
	Title  string
	Author string
	Videos []VideoInfo
}

var (
	// listCmd represents the list command
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "Print metadata of the desired playlist",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return checkOutputFormat()
		},
		Run: func(_ *cobra.Command, args []string) {
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

			exitOnError(writeOutput(os.Stdout, &playlistInfo, map[string]outputWriter{
				outputFormatPlain: func(w io.Writer) {
					writePlaylistOutput(w, &playlistInfo)
				},
				outputVideoIds: func(w io.Writer) {
					writePlaylistVideoIdsOutput(w, &playlistInfo)
				},
			}))
		},
	}
)

func writePlaylistVideoIdsOutput(w io.Writer, info *PlaylistInfo) {
	var ids []string
	for _, v := range info.Videos {
		ids = append(ids, v.ID)
	}
	fmt.Println(strings.Join(ids, " "))
}

func writePlaylistOutput(w io.Writer, info *PlaylistInfo) {
	fmt.Println("Title:      ", info.Title)
	fmt.Println("Author:     ", info.Author)
	fmt.Println("# Videos:   ", len(info.Videos))
	fmt.Println()

	table := tablewriter.NewWriter(w)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"ID", "Author", "Title", "Duration"})

	for _, vid := range info.Videos {
		table.Append([]string{
			vid.ID,
			vid.Author,
			vid.Title,
			fmt.Sprintf("%v", vid.Duration),
		})
	}
	table.Render()
}

func init() {
	rootCmd.AddCommand(listCmd)
	addFormatFlag(listCmd.Flags())
}
