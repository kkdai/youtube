package youtube

import (
	"context"
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
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

//SetLogOutput :Set logger writer
func SetLogOutput(w io.Writer) {
	log.SetOutput(w)
}

type stream struct {
	Quality string
	Type    string
	URL     string
	Itag    int
	Title   string
	Author  string
}

// Youtube implements the downloader to download youtube videos.
type Youtube struct {
	DebugMode         bool
	StreamList        []stream
	VideoID           string
	videoInfo         string
	DownloadPercent   chan int64
	Socks5Proxy       string
	contentLength     float64
	totalWrittenBytes float64
	downloadLevel     float64
}

//NewYoutube :Initialize youtube package object
func NewYoutube(debug bool) *Youtube {
	return &Youtube{DebugMode: debug, DownloadPercent: make(chan int64, 100)}
}

func NewYoutubeWithSocks5Proxy(debug bool, socks5Proxy string) *Youtube {
	return &Youtube{DebugMode: debug, DownloadPercent: make(chan int64, 100), Socks5Proxy: socks5Proxy}
}

//DecodeURL : Decode youtube URL to retrieval video information.
func (y *Youtube) DecodeURL(url string) error {
	err := y.findVideoID(url)
	if err != nil {
		return fmt.Errorf("findVideoID error=%s", err)
	}

	err = y.getVideoInfo()
	if err != nil {
		return fmt.Errorf("getVideoInfo error=%s", err)
	}

	err = y.parseVideoInfo()
	if err != nil {
		return fmt.Errorf("parse video info failed, err=%s", err)
	}

	return nil
}

//StartDownload : Starting download video to specific address.
func (y *Youtube) StartDownload(destFile string) error {
	//download highest resolution on [0]
	err := errors.New("Empty stream list")
	for _, stream := range y.StreamList {
		streamURL := stream.URL
		y.log(fmt.Sprintln("Download url=", streamURL))

		y.log(fmt.Sprintln("Download to file=", destFile))
		err = y.videoDLWorker(destFile, streamURL)
		if err != nil {
			y.log(fmt.Sprintln("Err: ", err.Error()))
			return err
		}
	}
	return err
}

//StartDownloadWithQuality : Starting download video with specific quality.
func (y *Youtube) StartDownloadWithQuality(destFile string, quality string) error {
	//download highest resolution on [0]
	err := errors.New("Empty stream list")
	for _, stream := range y.StreamList {
		if strings.Compare(stream.Quality, quality) == 0 {
			streamURL := stream.URL
			y.log(fmt.Sprintln("Download url=", streamURL))
			y.log(fmt.Sprintln("Download to file=", destFile))
			err = y.videoDLWorker(destFile, streamURL)
			if err == nil {
				break
			}
		}
	}

	if err != nil {
		return y.StartDownload(destFile)
	}
	return err
}

//StartDownloadFile : Starting download video on my download.
func (y *Youtube) StartDownloadFile() error {
	//download highest resolution on [0]
	err := errors.New("Empty stream list")
	for _, stream := range y.StreamList {
		streamURL := stream.URL
		y.log(fmt.Sprintln("Download url=", streamURL))

		// Find out what the file name should be.
		fileName := SanitizeFilename(stream.Title)
		fileName += pickIdealFileExtension(stream.Type)

		usr, _ := user.Current()
		destFile := filepath.Join(filepath.Join(usr.HomeDir, "Movies", "youtubedr"), fileName)
		y.log(fmt.Sprintln("Download to file=", destFile))

		err = y.videoDLWorker(destFile, streamURL)
		if err == nil {
			return nil
		}
	}
	return err
}

//  StartDownloadWithItag : Starting download video according to given itag value
func (y *Youtube) StartDownloadWithItag(destFile string, itag int) error {
	err := errors.New("Empty stream list")
	itagFound := false
	for _, stream := range y.StreamList {
		if stream.Itag == itag {
			itagFound = true
			streamURL := stream.URL
			y.log(fmt.Sprintln("Download url=", streamURL))

			if destFile == "" {
				// Find out what the file name should be.
				fileName := SanitizeFilename(stream.Title)
				fileName += pickIdealFileExtension(stream.Type)

				usr, _ := user.Current()
				destFile = filepath.Join(filepath.Join(usr.HomeDir, "Movies", "youtubedr"), fileName)
			}

			y.log(fmt.Sprintln("Download to file=", destFile))
			err = y.videoDLWorker(destFile, streamURL)
			if err == nil {
				return nil
			}
		}
	}

	if !itagFound {
		return errors.New("Invalid itag value, please specify correct value.")
	}
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
		err = fmt.Errorf("no response status found in the server's answer")
		return err
	}
	if status[0] == "fail" {
		reason, ok := answer["reason"]
		if ok {
			err = fmt.Errorf("'fail' response status found in the server's answer, reason: '%s'", reason[0])
		} else {
			err = errors.New("'fail' response status found in the server's answer, no reason given")
		}
		return err
	}
	if status[0] != "ok" {
		err = fmt.Errorf("non-success response status found in the server's answer (status: '%s')", status)
		return err
	}

	// read the streams map
	streamMap, ok := answer["player_response"]
	if !ok {
		err = errors.New("no stream map found in the server's answer")
		return err
	}

	// Get video title and author.
	title, author := getVideoTitleAuthor(answer)

	var prData PlayerResponseData
	if err := json.Unmarshal([]byte(streamMap[0]), &prData); err != nil {
		fmt.Println(err)
		panic("Player response json data has changed.")
	}

	// Get video download link
	if prData.PlayabilityStatus.Status == "UNPLAYABLE" {
		//Cannot playback on embedded video screen, could not download.
		return errors.New(fmt.Sprint("Cannot playback and download, reason:", prData.PlayabilityStatus.Reason))
	}

	var streams []stream
	for streamPos, streamRaw := range prData.StreamingData.Formats {
		if streamRaw.MimeType == "" {
			y.log(fmt.Sprintf("An error occurred while decoding one of the video's stream's information: stream %d.\n", streamPos))
			continue
		}
		streamUrl := streamRaw.URL
		if streamUrl == "" {
			cipher := streamRaw.Cipher
			decipheredUrl, err := y.decipher(cipher)
			if err != nil {
				return err
			}
			streamUrl = decipheredUrl
		}

		streams = append(streams, stream{
			Quality: streamRaw.Quality,
			Type:    streamRaw.MimeType,
			URL:     streamUrl,
			Itag:    streamRaw.Itag,

			Title:  title,
			Author: author,
		})
		y.log(fmt.Sprintf("Title: %s Author: %s Stream found: quality '%s', format '%s'", title, author, streamRaw.Quality, streamRaw.MimeType))
	}

	y.StreamList = streams
	if len(y.StreamList) == 0 {
		return errors.New("no stream list found in the server's answer")
	}
	return nil
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
	httpClient := &http.Client{Transport: httpTransport}

	if len(y.Socks5Proxy) == 0 {
		y.log("Using http without proxy.")
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

func (y *Youtube) getVideoInfo() error {
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
	defer resp.Body.Close()
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
		return errors.New("invalid characters in video id")
	}
	if len(videoID) < 10 {
		return errors.New("the video id must be at least 10 characters long")
	}
	return nil
}

func (y *Youtube) Write(p []byte) (n int, err error) {
	n = len(p)
	y.totalWrittenBytes = y.totalWrittenBytes + float64(n)
	currentPercent := ((y.totalWrittenBytes / y.contentLength) * 100)
	if (y.downloadLevel <= currentPercent) && (y.downloadLevel < 100) {
		y.downloadLevel++
		y.DownloadPercent <- int64(y.downloadLevel)
	}
	return
}
func (y *Youtube) videoDLWorker(destFile string, target string) error {

	httpClient, err := y.getHTTPClient()
	if err != nil {
		return err
	}

	resp, err := httpClient.Get(target)
	if err != nil {
		y.log(fmt.Sprintf("Http.Get\nerror: %s\ntarget: %s\n", err, target))
		return err
	}
	defer resp.Body.Close()
	y.contentLength = float64(resp.ContentLength)

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
	_, err = io.Copy(mw, resp.Body)
	if err != nil {
		y.log(fmt.Sprintln("download video err=", err))
		return err
	}
	return nil
}

func (y *Youtube) log(logText string) {
	if y.DebugMode {
		log.Println(logText)
	}
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
