package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (

	// listCmd represents the list command
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "Print metadata of the desired playlist",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			playlist, err := getDownloader().GetPlaylist(args[0])
			exitOnError(err)

			fmt.Println("Title:      ", playlist.Title)
			fmt.Println("Author:     ", playlist.Author)
			fmt.Println("# Videos:   ", len(playlist.Videos))

			fmt.Println()

			table := tablewriter.NewWriter(os.Stdout)
			table.SetAutoWrapText(false)
			table.SetHeader([]string{"ID", "Author", "Title", "Duration"})

			for _, vid := range playlist.Videos {

				table.Append([]string{
					vid.ID,
					vid.Author,
					vid.Title,
					fmt.Sprintf("%v", vid.Duration),
				})
			}
			table.Render()
		},
	}
)

func init() {
	infoCmd.Flags()
	rootCmd.AddCommand(listCmd)
}
