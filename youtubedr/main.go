package main

import (
	"flag"
	"fmt"
	"log"
	"os/user"
	"path/filepath"

	. "github.com/kkdai/youtube"
)

const usageString string = `Usage: youtubedr [OPTION] [URL]
Download a video from youtube.
Example: youtubedr -o "Campaign Diary".mp4 https://www.youtube.com/watch\?v\=XbNghLqsVwU
`

func main() {
	flag.Usage = func() {
		fmt.Println(usageString)
		flag.PrintDefaults()
	}
	usr, _ := user.Current()
	var outputFile string
	flag.StringVar(&outputFile, "o", "dl.mp4", "The output file")
	var outputDir string
	flag.StringVar(&outputDir, "d",
		filepath.Join(usr.HomeDir, "Movies", "yotubedr"),
		"The output directory.")
	flag.Parse()
	log.Println(flag.Args())
	log.Println("download to dir=", outputDir)
	y := NewYoutube(true)
	arg := flag.Arg(0)
	if err := y.DecodeURL(arg); err != nil {
		fmt.Println("err:", err)
		return
	}
	if err := y.StartDownload(filepath.Join(outputDir, outputFile)); err != nil {
		fmt.Println("err:", err)
	}
}
