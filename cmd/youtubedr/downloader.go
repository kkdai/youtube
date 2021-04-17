package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/kkdai/youtube/v2"
	ytdl "github.com/kkdai/youtube/v2/downloader"
	"github.com/spf13/pflag"
	"golang.org/x/net/http/httpproxy"
)

var (
	insecureSkipVerify bool   // skip TLS server validation
	outputQuality      string // itag number or quality string
	mimetype           string // mimetype
	downloader         *ytdl.Downloader
)

func addQualityFlag(flagSet *pflag.FlagSet) {
	flagSet.StringVarP(&outputQuality, "quality", "q", "", "The itag number or quality label (hd720, medium)")
}

func addMimeTypeFlag(flagSet *pflag.FlagSet) {
	flagSet.StringVarP(&mimetype, "mimetype", "m", "mp4", "Mime-Type to filter (mp4, webm, av01, avc1)")
}

func getDownloader() *ytdl.Downloader {
	if downloader != nil {
		return downloader
	}

	proxyFunc := httpproxy.FromEnvironment().ProxyFunc()
	httpTransport := &http.Transport{
		// Proxy: http.ProxyFromEnvironment() does not work. Why?
		Proxy: func(r *http.Request) (uri *url.URL, err error) {
			return proxyFunc(r.URL)
		},
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
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
	downloader.Client.Debug = verbose
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

	formats := video.Formats
	if mimetype != "" {
		formats = formats.Type(mimetype)
	}
	if len(formats) == 0 {
		return nil, nil, errors.New("no formats found")
	}

	var format *youtube.Format
	switch {
	case itag > 0:
		format = formats.FindByItag(itag)
		if format == nil {
			return nil, nil, fmt.Errorf("unable to find format with itag %d", itag)
		}

	case outputQuality != "":
		format = formats.FindByQuality(outputQuality)
		if format == nil {
			return nil, nil, fmt.Errorf("unable to find format with quality %s", outputQuality)
		}

	default:
		// select the first format
		formats.Sort()
		format = &formats[0]
	}

	return video, format, nil
}
