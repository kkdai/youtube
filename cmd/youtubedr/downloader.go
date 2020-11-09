package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/kkdai/youtube/v2"
	ytdl "github.com/kkdai/youtube/v2/downloader"
	"github.com/spf13/pflag"
)

var (
	insecureSkipVerify bool   // skip TLS server validation
	outputQuality      string // itag number or quality string
	downloader         *ytdl.Downloader
)

func addQualityFlag(flagSet *pflag.FlagSet) {
	flagSet.StringVarP(&outputQuality, "quality", "q", "", "The itag number or quality label (hd720, medium)")
}

func getDownloader() *ytdl.Downloader {
	if downloader != nil {
		return downloader
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

	downloader = &ytdl.Downloader{
		OutputDir: outputDir,
	}
	downloader.HTTPClient = &http.Client{Transport: httpTransport}

	return downloader
}

func getVideoWithFormat(id string) (*youtube.Video, *youtube.Format, error) {
	dl := getDownloader()
	itag, _ := strconv.Atoi(outputQuality)

	video, err := dl.GetVideo(id)
	if err != nil {
		return nil, nil, err
	}

	if len(video.Formats) == 0 {
		return nil, nil, errors.New("no formats found")
	}

	var format *youtube.Format
	switch {
	case itag > 0:
		format = video.Formats.FindByItag(itag)
		if format == nil {
			return nil, nil, fmt.Errorf("unable to find format with itag %d", itag)
		}

	case outputQuality != "":
		format = video.Formats.FindByQuality(outputQuality)
		if format == nil {
			return nil, nil, fmt.Errorf("unable to find format with quality %s", outputQuality)
		}

	default:
		// select the first format
		format = &video.Formats[0]
	}

	return video, format, nil
}
