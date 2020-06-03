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
	flag.StringVar(&outputFile, "o", "", "The output file")
	var outputDir string
	flag.StringVar(&outputDir, "d",
		filepath.Join(usr.HomeDir, "Movies", "youtubedr"),
		"The output directory.")
	var outputQuality string
	flag.StringVar(&outputQuality, "q", "", "The output file quality (hd720, medium)")

	var socks5Proxy string
	flag.StringVar(&socks5Proxy, "p", "", "The Socks 5 proxy, e.g. 10.10.10.10:7878")

	var itag int
	flag.IntVar(&itag, "i", 0, "Specify itag number, e.g. 13, 17")

	flag.Parse()
	log.Println(flag.Args())
	log.Println("download to dir=", outputDir)
	y := NewYoutubeWithSocks5Proxy(true, socks5Proxy)
	arg := flag.Arg(0)
	if err := y.DecodeURL(arg); err != nil {
		fmt.Println("err:", err)
		return
	}

	err := y.StartDownload(outputDir, outputFile, outputQuality, itag)
	if err != nil {
		fmt.Println("err:", err)
	}
}
