package main

import (
	"fmt"
	"os"
)

func main() {
	exitOnError(rootCmd.Execute())
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func exitOnErrors(errors []error) {
	for _, err := range errors {
		fmt.Fprintln(os.Stderr, err)
	}
	if len(errors) != 0 {
		os.Exit(1)
	}
}
