package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/kkdai/youtube"
)

const usageString string = `Usage: youtubedr [OPTION] [URL]
Download a video from youtube.
Example: youtubedr -o "Campaign Diary".mp4 https://www.youtube.com/watch\?v\=XbNghLqsVwU
`

var (
	outputFile         string
	outputDir          string
	outputQuality      string
	socks5Proxy        string
	itag               int
	info               bool
	insecureSkipVerify bool
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		os.Exit(1)
	}
}

func run() error {
	flag.Usage = func() {
		fmt.Println(usageString)
		flag.PrintDefaults()
	}
	usr, _ := user.Current()
	flag.StringVar(&outputFile, "o", "", "The output file")
	flag.StringVar(&outputDir, "d",
		filepath.Join(usr.HomeDir, "Movies", "youtubedr"),
		"The output directory.")
	flag.StringVar(&outputQuality, "q", "", "The output file quality (hd720, medium)")
	flag.StringVar(&socks5Proxy, "p", "", "The Socks 5 proxy, e.g. 10.10.10.10:7878")
	flag.IntVar(&itag, "i", 0, "Specify itag number, e.g. 13, 17")
	flag.BoolVar(&info, "info", false, "show info of video")
	flag.BoolVar(&insecureSkipVerify, "insecure-skip-tls-verify", false, "skip server certificate verification")

	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.PrintDefaults()
		return nil
	}

	log.Println("download to dir=", outputDir)
	y := youtube.NewYoutubeWithSocks5Proxy(true, socks5Proxy, insecureSkipVerify)
	if len(y.Socks5Proxy) == 0 {
		log.Println("Using http without proxy.")
	}
	if y.InsecureSkipVerify {
		log.Println("Skip server certificate verification")
	}
	arg := flag.Arg(0)
	if err := y.DecodeURL(arg); err != nil {
		return err
	}

	if info {
		info := y.GetStreamInfo()
		if info == nil {
			fmt.Println("-----no available stream-----")
			return nil
		}
		fmt.Printf("Title: %s\n", info.Title)
		fmt.Printf("Author: %s\n", info.Author)
		fmt.Println("-----available streams-----")
		for _, itag := range info.Streams {
			fmt.Printf("itag: %3d , quality: %6s , type: %10s\n", itag.ItagNo, itag.Quality, itag.MimeType)
		}
		return nil
	}

	if outputQuality == "hd1080" {
		fmt.Println("check ffmpeg is installed....")
		ffmpegVersionCmd := exec.Command("ffmpeg", "-version")
		if err := ffmpegVersionCmd.Run(); err != nil {
			return fmt.Errorf("please check ffmpeg is installed correctly, err: %w", err)
		}
		return y.StartDownloadWithHighQuality(outputDir, outputFile, outputQuality)
	}
	return y.StartDownload(outputDir, outputFile, outputQuality, itag)
}
