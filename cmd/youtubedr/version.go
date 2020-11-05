package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// set through ldflags
	version string
	commit  string
	date    string
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:    ", version)
		fmt.Println("Commit:     ", commit)
		fmt.Println("Date:       ", date)
		fmt.Println("Go Version: ", runtime.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
