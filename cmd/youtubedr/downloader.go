package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/spf13/pflag"
	"golang.org/x/net/http/httpproxy"

	"github.com/kkdai/youtube/v2"
	ytdl "github.com/kkdai/youtube/v2/downloader"
)

var (
	insecureSkipVerify bool   // skip TLS server validation
	outputQuality      string // itag number or quality string
	mimetype           string // mimetype
	downloader         *ytdl.Downloader
)

func addQualityFlag(flagSet *pflag.FlagSet) {
	flagSet.StringVarP(&outputQuality, "quality", "q", "medium", "The itag number or quality label (hd720, medium)")
}

func addMimeTypeFlag(flagSet *pflag.FlagSet) {
	flagSet.StringVarP(&mimetype, "mimetype", "m", "mp4", "Mime-Type to filter (mp4, webm, av01, avc1) - applicable if --quality used is quality label")
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

	youtube.SetLogLevel(logLevel)

	if insecureSkipVerify {
		youtube.Logger.Info("Skip server certificate verification")
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
	itag, _ := strconv.Atoi(outputQuality)
	switch {
	case itag > 0:
		// When an itag is specified, do not filter format with mime-type
		format = video.Formats.FindByItag(itag)
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
