package main

import (
	"crypto/tls"
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
	mimetype           string
	language           string
	downloader         *ytdl.Downloader
)

func addVideoSelectionFlags(flagSet *pflag.FlagSet) {
	flagSet.StringVarP(&outputQuality, "quality", "q", "medium", "The itag number or quality label (hd720, medium)")
	flagSet.StringVarP(&mimetype, "mimetype", "m", "", "Mime-Type to filter (mp4, webm, av01, avc1) - applicable if --quality used is quality label")
	flagSet.StringVarP(&language, "language", "l", "", "Language to filter")
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

func getVideoWithFormat(videoID string) (*youtube.Video, *youtube.Format, error) {
	dl := getDownloader()
	video, err := dl.GetVideo(videoID)
	if err != nil {
		return nil, nil, err
	}

	itag, _ := strconv.Atoi(outputQuality)
	formats := video.Formats

	if language != "" {
		formats = formats.Language(language)
	}
	if mimetype != "" {
		formats = formats.Type(mimetype)
	}
	if outputQuality != "" {
		formats = formats.Quality(outputQuality)
	}
	if itag > 0 {
		formats = formats.Itag(itag)
	}
	if formats == nil {
		return nil, nil, fmt.Errorf("unable to find the specified format")
	}

	formats.Sort()

	// select the first format
	return video, &formats[0], nil
}
