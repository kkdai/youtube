package youtube

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/kkdai/youtube/pkg/decipher"

	"github.com/google/uuid"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"

	"golang.org/x/net/proxy"
)

//SetLogOutput :Set logger writer
func SetLogOutput(w io.Writer) {
	log.SetOutput(w)
}

// Youtube implements the downloader to download youtube videos.
type Youtube struct {
	DebugMode         bool
	StreamList        []Stream
	VideoID           string
	videoInfo         string
	DownloadPercent   chan int64
	Socks5Proxy       string
	contentLength     float64
	totalWrittenBytes float64
	downloadLevel     float64
	Title             string
	Author            string

	// Disable SSL certificate verification
	InsecureSkipVerify bool
}

//NewYoutube :Initialize youtube package object.
func NewYoutube(debug bool, insecureSkipVerify bool) *Youtube {
	return &Youtube{DebugMode: debug, DownloadPercent: make(chan int64, 100), InsecureSkipVerify: insecureSkipVerify}
}

func NewYoutubeWithSocks5Proxy(debug bool, socks5Proxy string, insecureSkipVerify bool) *Youtube {
	return &Youtube{DebugMode: debug, DownloadPercent: make(chan int64, 100), Socks5Proxy: socks5Proxy, InsecureSkipVerify: insecureSkipVerify}
}

//DecodeURL : Decode youtube URL to retrieval video information.
func (y *Youtube) DecodeURL(url string) error {
	err := y.findVideoID(url)
	if err != nil {
		return fmt.Errorf("findVideoID failed: %w", err)
	}

	err = y.getVideoInfo()
	if err != nil {
		return fmt.Errorf("getVideoInfo failed: %w", err)
	}

	err = y.parseVideoInfo()
	if err != nil {
		return fmt.Errorf("parse video info failed: %w", err)
	}

	return nil
}

//StartDownload : Starting download video by arguments.
func (y *Youtube) StartDownload(outputDir, outputFile, quality string, itagNo int) error {
	if len(y.StreamList) == 0 {
		return ErrEmptyStreamList
	}

	//download highest resolution on [0] by default
	index := 0
	switch {
	case itagNo != 0:
		itagFound := false
		for i, stream := range y.StreamList {
			if stream.ItagNo == itagNo {
				itagFound = true
				index = i
				break
			}
		}
		if !itagFound {
			return ErrItagNotFound
		}
	case quality != "":
		for i, stream := range y.StreamList {
			if strings.Compare(stream.Quality, quality) == 0 {
				index = i
				break
			}
		}
	}
	stream := y.StreamList[index]

	if outputDir == "" {
		usr, _ := user.Current()
		outputDir = filepath.Join(usr.HomeDir, "Movies", "youtubedr")
	}

	outputFile = SanitizeFilename(outputFile)
	if outputFile == "" {
		outputFile = SanitizeFilename(y.Title)
		outputFile += pickIdealFileExtension(stream.MimeType)
	}
	destFile := filepath.Join(outputDir, outputFile)
	streamURL, err := y.getStreamUrl(stream)
	if err != nil {
		return err
	}
	y.log(fmt.Sprintln("Download url=", streamURL))
	y.log(fmt.Sprintln("Download to file=", destFile))
	return y.videoDLWorker(destFile, streamURL)
}

func (y *Youtube) getStreamUrl(stream Stream) (string, error) {
	streamURL := stream.URL
	if streamURL == "" {
		cipher := stream.Cipher
		if cipher == "" {
			return "", ErrCipherNotFound
		}
		client, err := y.getHTTPClient()
		if err != nil {
			return "", fmt.Errorf("get http client failed: %w", err)
		}
		decipher := decipher.NewDecipher(client)
		decipherUrl, err := decipher.Url(y.VideoID, cipher)
		if err != nil {
			return "", err
		}
		streamURL = decipherUrl
	}
	return streamURL, nil
}

//StartDownloadWithHighQuality : Starting downloading video with high quality (>720p).
func (y *Youtube) StartDownloadWithHighQuality(outputDir string, outputFile string, quality string) error {
	if len(y.StreamList) == 0 {
		return ErrEmptyStreamList
	}

	qualityitagMap := map[string]struct {
		videoItag int
		audioItag int
	}{
		"hd1080": {137, 140},
	}
	videoItag := qualityitagMap[quality].videoItag
	audioItag := qualityitagMap[quality].audioItag

	var videoStream, audioStream Stream

	for _, stream := range y.StreamList {
		switch stream.ItagNo {
		case videoItag:
			videoStream = stream
		case audioItag:
			audioStream = stream
		}
	}

	if videoStream.ItagNo == 0 {
		return errors.New("no Stream video/mp4 for hd1080 found")
	}
	if audioStream.ItagNo == 0 {
		return errors.New("no Stream audio/mp4 for hd1080 found")
	}

	if outputDir == "" {
		usr, _ := user.Current()
		outputDir = filepath.Join(usr.HomeDir, "Movies", "youtubedr")
	}

	outputFile = SanitizeFilename(outputFile)
	stream := videoStream
	if outputFile == "" {
		outputFile = SanitizeFilename(y.Title)
		outputFile += pickIdealFileExtension(stream.MimeType)
	}
	uid := uuid.New()
	tempFileName := "temp_" + uid.String()
	videoFile := filepath.Join(outputDir, tempFileName+".m4v")
	audioFile := filepath.Join(outputDir, tempFileName+".m4a")
	defer func() {
		if err := os.Remove(videoFile); err != nil {
			y.log(fmt.Sprintf("err to remove file: %s", err))
		}
		if err := os.Remove(audioFile); err != nil {
			y.log(fmt.Sprintf("err to remove file: %s", err))
		}
	}()
	var err error
	videoStreamUrl, err := y.getStreamUrl(videoStream)
	if err != nil {
		return err
	}
	y.log(fmt.Sprintln("Download url=", videoStreamUrl))
	y.log("Downloading video file...")
	err = y.videoDLWorker(videoFile, videoStreamUrl)
	if err != nil {
		return err
	}
	audioStreamUrl, err := y.getStreamUrl(audioStream)
	if err != nil {
		return err
	}
	y.log(fmt.Sprintln("Download url=", audioStreamUrl))
	y.log("Downloading audio file...")
	err = y.videoDLWorker(audioFile, audioStreamUrl)
	if err != nil {
		return err
	}

	destFile := filepath.Join(outputDir, outputFile)
	y.log(fmt.Sprintln("Download to file=", destFile))
	ffmpegVersionCmd := exec.Command("ffmpeg", "-y", "-i", videoFile, "-i", audioFile, "-strict", "-2", "-shortest", destFile, "-loglevel", "warning")
	ffmpegVersionCmd.Stderr = os.Stderr
	ffmpegVersionCmd.Stdout = os.Stdout
	y.log("merging video and audio.....")
	if err := ffmpegVersionCmd.Run(); err != nil {
		fmt.Println("err:", err)
		os.Exit(1)
	}
	y.log("Done")
	return err
}

func pickIdealFileExtension(mediaType string) string {
	defaultExtension := ".mov"

	mediaType, _, err := mime.ParseMediaType(mediaType)
	if err != nil {
		return defaultExtension
	}

	// Rely on hardcoded canonical mime types, as the ones provided by Go aren't exhaustive [1].
	// This seems to be a recurring problem for youtube downloaders, see [2].
	// The implementation is based on mozilla's list [3], IANA [4] and Youtube's support [5].
	// [1] https://github.com/golang/go/blob/ed7888aea6021e25b0ea58bcad3f26da2b139432/src/mime/type.go#L60
	// [2] https://github.com/ZiTAL/youtube-dl/blob/master/mime.types
	// [3] https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
	// [4] https://www.iana.org/assignments/media-types/media-types.xhtml#video
	// [5] https://support.google.com/youtube/troubleshooter/2888402?hl=en
	canonicals := map[string]string{
		"video/quicktime":  ".mov",
		"video/x-msvideo":  ".avi",
		"video/x-matroska": ".mkv",
		"video/mpeg":       ".mpeg",
		"video/webm":       ".webm",
		"video/3gpp2":      ".3g2",
		"video/x-flv":      ".flv",
		"video/3gpp":       ".3gp",
		"video/mp4":        ".mp4",
		"video/ogg":        ".ogv",
		"video/mp2t":       ".ts",
	}

	if extension, ok := canonicals[mediaType]; ok {
		return extension
	}

	// Our last resort is to ask the operating system, but these give multiple results and are rarely canonical.
	extensions, err := mime.ExtensionsByType(mediaType)
	if err != nil || extensions == nil {
		return defaultExtension
	}

	return extensions[0]
}

func SanitizeFilename(fileName string) string {
	// Characters not allowed on mac
	//	:/
	// Characters not allowed on linux
	//	/
	// Characters not allowed on windows
	//	<>:"/\|?*

	// Ref https://docs.microsoft.com/en-us/windows/win32/fileio/naming-a-file#naming-conventions

	fileName = regexp.MustCompile(`[:/<>\:"\\|?*]`).ReplaceAllString(fileName, "")
	fileName = regexp.MustCompile(`\s+`).ReplaceAllString(fileName, " ")

	return fileName
}

func (y *Youtube) parseVideoInfo() error {
	answer, err := url.ParseQuery(y.videoInfo)
	if err != nil {
		return err
	}

	status, ok := answer["status"]
	if !ok {
		return fmt.Errorf("no response status found in the server's answer")
	}
	if status[0] == "fail" {
		reason, ok := answer["reason"]
		if ok {
			return fmt.Errorf("'fail' response status found in the server's answer, reason: '%s'", reason[0])
		}
		return errors.New("'fail' response status found in the server's answer, no reason given")
	}
	if status[0] != "ok" {
		return fmt.Errorf("non-success response status found in the server's answer (status: '%s')", status)
	}

	// read the streams map
	streamMap, ok := answer["player_response"]
	if !ok {
		return errors.New("no Stream map found in the server's answer")
	}

	var prData PlayerResponseData
	if err := json.Unmarshal([]byte(streamMap[0]), &prData); err != nil {
		fmt.Println(err)
		panic("Player response json data has changed.")
	}

	// Check if video is downloadable
	if prData.PlayabilityStatus.Status != "OK" {
		return fmt.Errorf("cannot playback and download, status: %s, reason: %s", prData.PlayabilityStatus.Status, prData.PlayabilityStatus.Reason)
	}

	// Get video title and author.
	y.Title, y.Author = getVideoTitleAuthor(answer)

	// Get video download link
	streams, err := y.getStreams(prData)
	if err != nil {
		return err
	}

	y.StreamList = streams
	if len(y.StreamList) == 0 {
		return errors.New("no Stream list found in the server's answer")
	}

	return nil
}

func (y Youtube) getStreams(prData PlayerResponseData) ([]Stream, error) {
	size := len(prData.StreamingData.Formats) + len(prData.StreamingData.AdaptiveFormats)
	streams := make([]Stream, 0, size)

	filterFormat := func(stream Stream) {
		if stream.MimeType == "" {
			y.log(fmt.Sprintf("An error occurred while decoding one of the video's Stream's information: Stream %+v.\n", stream))
			return
		}
		streams = append(streams, stream)
	}

	for _, format := range prData.StreamingData.Formats {
		filterFormat(format.Stream)
	}
	for _, format := range prData.StreamingData.AdaptiveFormats {
		filterFormat(format.Stream)
	}
	return streams, nil
}

func (y *Youtube) getHTTPClient() (*http.Client, error) {
	// setup a http client
	httpTransport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	if y.InsecureSkipVerify {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // nolint
	}
	httpClient := &http.Client{Transport: httpTransport}

	if len(y.Socks5Proxy) == 0 {
		return httpClient, nil
	}

	dialer, err := proxy.SOCKS5("tcp", y.Socks5Proxy, nil, proxy.Direct)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		return nil, err
	}
	// set our socks5 as the dialer
	dc := dialer.(interface {
		DialContext(ctx context.Context, network, addr string) (net.Conn, error)
	})
	httpTransport.DialContext = dc.DialContext

	y.log(fmt.Sprintf("Using http with proxy %s.", y.Socks5Proxy))

	return httpClient, nil
}

func (y *Youtube) getVideoInfo() (err error) {
	eurl := "https://youtube.googleapis.com/v/" + y.VideoID
	url := "https://youtube.com/get_video_info?video_id=" + y.VideoID + "&eurl=" + eurl
	y.log(fmt.Sprintf("url: %s", url))

	httpClient, err := y.getHTTPClient()
	if err != nil {
		return err
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer y.Close(resp.Body, "getVideoInfo")
	if resp.StatusCode != 200 {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	y.videoInfo = string(body)
	return nil
}

func (y Youtube) Close(r io.ReadCloser, op string) {
	_, err := io.Copy(ioutil.Discard, r)
	if err != nil && err.Error() != ErrReadOnClosedResBody.Error() {
		y.log(fmt.Sprintf("failed to exhaust reader: %s in %s", err, op))
	}
	err = r.Close()
	if err != nil {
		y.log(fmt.Sprintf("response close err %s in %s", err, op))
	}
}

func (y *Youtube) findVideoID(url string) error {
	videoID := url
	if strings.Contains(videoID, "youtu") || strings.ContainsAny(videoID, "\"?&/<%=") {
		reList := []*regexp.Regexp{
			regexp.MustCompile(`(?:v|embed|watch\?v)(?:=|/)([^"&?/=%]{11})`),
			regexp.MustCompile(`(?:=|/)([^"&?/=%]{11})`),
			regexp.MustCompile(`([^"&?/=%]{11})`),
		}
		for _, re := range reList {
			if isMatch := re.MatchString(videoID); isMatch {
				subs := re.FindStringSubmatch(videoID)
				videoID = subs[1]
			}
		}
	}
	y.log(fmt.Sprintf("Found video id: '%s'", videoID))
	y.VideoID = videoID
	if strings.ContainsAny(videoID, "?&/<%=") {
		return ErrInvalidCharactersInVideoId
	}
	if len(videoID) < 10 {
		return ErrVideoIdMinLength
	}
	return nil
}

func (y *Youtube) Write(p []byte) (n int, err error) {
	n = len(p)
	y.totalWrittenBytes = y.totalWrittenBytes + float64(n)
	currentPercent := (y.totalWrittenBytes / y.contentLength) * 100
	if (y.downloadLevel <= currentPercent) && (y.downloadLevel < 100) {
		y.downloadLevel++
		y.DownloadPercent <- int64(y.downloadLevel)
	}
	return
}
func (y *Youtube) videoDLWorker(destFile string, target string) (err error) {
	httpClient, err := y.getHTTPClient()
	if err != nil {
		return err
	}

	resp, err := httpClient.Get(target)
	if err != nil {
		y.log(fmt.Sprintf("Http.Get\nerror: %s\ntarget: %s\n", err, target))
		return err
	}
	defer y.Close(resp.Body, "videoDLWorker")
	y.contentLength = float64(resp.ContentLength)

	// create progress bar
	progress := mpb.New(mpb.WithWidth(64))
	bar := progress.AddBar(
		int64(y.contentLength),

		mpb.PrependDecorators(
			decor.CountersKibiByte("% .2f / % .2f"),
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_GO, 90),
			decor.Name(" ] "),
			decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
		),
	)
	reader := bar.ProxyReader(resp.Body)
	defer y.Close(reader, "progress bar")

	if resp.StatusCode != 200 {
		y.log(fmt.Sprintf("reading answer: non 200[code=%v] status code received: '%v'", resp.StatusCode, err))
		return errors.New("non 200 status code received")
	}
	err = os.MkdirAll(filepath.Dir(destFile), 0755)
	if err != nil {
		return err
	}
	out, err := os.Create(destFile)
	if err != nil {
		return err
	}
	mw := io.MultiWriter(out, y)
	_, err = io.Copy(mw, reader)
	if err != nil {
		y.log(fmt.Sprintln("download video err=", err))
		return err
	}
	progress.Wait()
	return nil
}

func (y *Youtube) log(logText string) {
	if y.DebugMode {
		log.Println(logText)
	}
}

func (y *Youtube) GetStreamInfo() *StreamInfo {
	if len(y.StreamList) == 0 {
		return nil
	}
	return &StreamInfo{Title: y.Title, Author: y.Author, Streams: y.StreamList}
}

func getVideoTitleAuthor(in url.Values) (string, string) {
	playResponse, ok := in["player_response"]
	if !ok {
		return "", ""
	}
	personMap := make(map[string]interface{})

	if err := json.Unmarshal([]byte(playResponse[0]), &personMap); err != nil {
		panic(err)
	}

	s := personMap["videoDetails"]
	myMap := s.(map[string]interface{})
	// fmt.Println("-->", myMap["title"], "oooo:", myMap["author"])
	if title, ok := myMap["title"]; ok {
		if author, ok := myMap["author"]; ok {
			return title.(string), author.(string)
		}
	}

	return "", ""
}
