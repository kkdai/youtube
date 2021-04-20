package main

import (
	"encoding/json"
	"encoding/xml"
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

type playlistOutputWriter func(PlaylistInfo, io.Writer) error

var (
	playlistOutputFunc    playlistOutputWriter
	playlistOutputWriters = map[string]playlistOutputWriter{
		"json": func(info PlaylistInfo, w io.Writer) error {
			encoder := json.NewEncoder(w)
			encoder.SetIndent("", "  ")
			return encoder.Encode(info)
		},
		"xml": func(info PlaylistInfo, w io.Writer) error {
			return xml.NewEncoder(w).Encode(info)
		},
		"plain": func(info PlaylistInfo, w io.Writer) error {
			fmt.Println("Title:      ", info.Title)
			fmt.Println("Author:     ", info.Author)
			fmt.Println("# Videos:   ", len(info.Videos))

			fmt.Println()

			table := tablewriter.NewWriter(os.Stdout)
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
			return nil
		},
	}
	// listCmd represents the list command
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "Print metadata of the desired playlist",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			playlistOutputFunc = playlistOutputWriters[outputFormat]
			if playlistOutputFunc == nil {
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
			exitOnError(playlistOutputFunc(playlistInfo, os.Stdout))
		},
	}
)

func init() {
	rootCmd.AddCommand(listCmd)
	addFormatFlag(listCmd.Flags())
}
