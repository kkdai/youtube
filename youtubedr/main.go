package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"time"

	ytdl "github.com/kkdai/youtube/downloader"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/net/proxy"
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

	httpTransport := &http.Transport{
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if socks5Proxy != "" {
		log.Println("Using SOCKS5 proxy", socks5Proxy)
		dialer, err := proxy.SOCKS5("tcp", socks5Proxy, nil, proxy.Direct)
		if err != nil {
			fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
			os.Exit(1)
		}

		// set our socks5 as the dialer
		dc := dialer.(interface {
			DialContext(ctx context.Context, network, addr string) (net.Conn, error)
		})
		httpTransport.DialContext = dc.DialContext
	} else {
		httpTransport.DialContext = (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext
	}

	if insecureSkipVerify {
		log.Println("Skip server certificate verification")
		httpTransport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	dl := ytdl.Downloader{
		OutputDir: outputDir,
	}
	dl.Client.HTTPClient = &http.Client{Transport: httpTransport}

	arg := flag.Arg(0)

	video, err := dl.Client.GetVideo(arg)
	if err != nil {
		return err
	}

	if info {
		fmt.Printf("Title:    %s\n", video.Title)
		fmt.Printf("Author:   %s\n", video.Author)
		fmt.Printf("Duration: %v\n", video.Duration)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoWrapText(false)
		table.SetHeader([]string{"itag", "quality", "MimeType"})

		for _, itag := range video.Streams {
			table.Append([]string{strconv.Itoa(itag.ItagNo), itag.Quality, itag.MimeType})
		}
		table.Render()
		return nil
	}

	fmt.Println("download to directory", outputDir)

	if outputQuality == "hd1080" {
		fmt.Println("check ffmpeg is installed....")
		ffmpegVersionCmd := exec.Command("ffmpeg", "-version")
		if err := ffmpegVersionCmd.Run(); err != nil {
			return fmt.Errorf("please check ffmpeg is installed correctly, err: %w", err)
		}
		return dl.DownloadWithHighQuality(context.Background(), outputFile, video, outputQuality)
	}

	return dl.Download(context.Background(), outputFile, video, &video.Streams[0])
}
