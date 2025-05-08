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
	isDownloader bool
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
			if !isDownloader {
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

				exitOnError(writeOutput(os.Stdout, &playlistInfo, func(w io.Writer) {
					writePlaylistOutput(w, &playlistInfo)
				}))
			} else {
				for i := 0; i < len(playlist.Videos); i++ {
					err := download(playlist.Videos[i].ID)
					if err != nil {
						fmt.Println(err.Error())
					}
					err = nil
				}
			}
		},
	}
)

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

	listCmd.Flags().BoolVarP(&isDownloader, "downloader", "d", false, "Download playlist")
	listCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "The output directory.")
	addFormatFlag(listCmd.Flags())
}
