package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/kkdai/youtube/v2"
	ytdl "github.com/kkdai/youtube/v2/downloader"
	"github.com/olekukonko/tablewriter"
)

const usageString string = `Usage: youtubedr [OPTION] [URL]
Download a video from youtube.
Example: youtubedr -o "Campaign Diary".mp4 https://www.youtube.com/watch\?v\=XbNghLqsVwU
`

var (
	outputFile         string
	outputDir          string
	outputQuality      string
	itag               int
	info               bool
	insecureSkipVerify bool
	onlyStreamUrl      bool
	showVersion        bool

	// set through ldflags
	version string
	commit  string
	date    string
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
		fmt.Println("\n" + `Use the HTTP_PROXY environment variable to set a HTTP or SOCSK5 proxy. The proxy type is determined by the URL scheme.
"http", "https", and "socks5" are supported. If the scheme is empty, "http" is assumed."`)
	}
	flag.StringVar(&outputFile, "o", "", "The output file")
	flag.StringVar(&outputDir, "d", ".", "The output directory.")
	flag.StringVar(&outputQuality, "q", "", "The output file quality (hd720, medium)")
	flag.IntVar(&itag, "i", 0, "Specify itag number, e.g. 13, 17")
	flag.BoolVar(&info, "info", false, "show info of video")
	flag.BoolVar(&showVersion, "version", false, "show version info and exit")
	flag.BoolVar(&insecureSkipVerify, "insecure-skip-tls-verify", false, "skip server certificate verification")
	flag.BoolVar(&onlyStreamUrl, "only-stream-url", false, "Only output the stream-url to desired video")

	flag.Parse()

	if onlyStreamUrl {
		log.SetOutput(ioutil.Discard)
	}

	if showVersion {
		fmt.Println("Version:    ", version)
		fmt.Println("Commit:     ", commit)
		fmt.Println("Date:       ", date)
		fmt.Println("Go Version: ", runtime.Version())
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		flag.Usage()
		return nil
	}

	httpTransport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
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
	dl.HTTPClient = &http.Client{Transport: httpTransport}

	arg := flag.Arg(0)

	video, err := dl.GetVideo(arg)
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

		for _, itag := range video.Formats {
			table.Append([]string{strconv.Itoa(itag.ItagNo), itag.Quality, itag.MimeType})
		}
		table.Render()
		return nil
	}

	log.Println("download to directory", outputDir)

	if len(video.Formats) == 0 {
		return errors.New("no formats found")
	}

	var format *youtube.Format
	if itag > 0 {
		format = video.Formats.FindByItag(itag)
		if format == nil {
			return fmt.Errorf("unable to find format with itag %d", itag)
		}
		outputQuality = format.Quality
	} else if outputQuality != "" {
		format = video.Formats.FindByQuality(outputQuality)
		if format == nil {
			return fmt.Errorf("unable to find format with quality %s", outputQuality)
		}
	} else {
		format = &video.Formats[0]
	}

	if onlyStreamUrl {
		url, _ := dl.GetStreamURL(video, format)
		fmt.Println(url)
		return nil
	}

	if outputQuality == "hd1080" {
		fmt.Println("check ffmpeg is installed....")
		ffmpegVersionCmd := exec.Command("ffmpeg", "-version")
		if err := ffmpegVersionCmd.Run(); err != nil {
			return fmt.Errorf("please check ffmpeg is installed correctly, err: %w", err)
		}
		return dl.DownloadWithHighQuality(context.Background(), outputFile, video, outputQuality)
	}
	return dl.Download(context.Background(), video, format, outputFile)
}
