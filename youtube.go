package youtube

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func NewYoutube() *Youtube {
	return new(Youtube)
}

type stream map[string]string

type Youtube struct {
	StreamList []stream
	videoID    string
	videoInfo  string
}

func (y *Youtube) DecodeURL(url string) error {
	err := y.findVideoId(url)
	if err != nil {
		return fmt.Errorf("findvideoID error=%s", err)
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

func (y *Youtube) StartDownload(dstDir string) {
	//download highest resolution on [0]
	targetStream := y.streamList[0]
	url := targetStream["url"] + "&signature=" + targetStream["sig"]
	log.Println("Download url=", url)

	targetFile := fmt.Sprintf("%s/%s.%s", dstDir, targetStream["title"], "mp4")
	//targetStream["title"], targetStream["author"])
	log.Println("Download to file=", targetFile)
	videoDLWorker(targetFile, url)
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
			err = errors.New(fmt.Sprint("'fail' response status found in the server's answer, no reason given"))
		}
		return err
	}
	if status[0] != "ok" {
		err = fmt.Errorf("non-success response status found in the server's answer (status: '%s')", status)
		return err
	}

	// read the streams map
	stream_map, ok := answer["url_encoded_fmt_stream_map"]
	if !ok {
		err = errors.New(fmt.Sprint("no stream map found in the server's answer"))
		return err
	}

	// read each stream

	streams_list := strings.Split(stream_map[0], ",")

	var streams []stream
	for stream_pos, stream_raw := range streams_list {
		stream_qry, err := url.ParseQuery(stream_raw)
		if err != nil {
			log.Println(fmt.Sprintf("An error occured while decoding one of the video's stream's information: stream %d: %s\n", stream_pos, err))
			continue
		}
		var sig string
		if _, exist := stream_qry["sig"]; exist {
			sig = stream_qry["sig"][0]
		}

		stream := stream{
			"quality": stream_qry["quality"][0],
			"type":    stream_qry["type"][0],
			"url":     stream_qry["url"][0],
			"sig":     sig,
			"title":   answer["title"][0],
			"author":  answer["author"][0],
		}
		streams = append(streams, stream)
		log.Printf("Stream found: quality '%s', format '%s'", stream_qry["quality"][0], stream_qry["type"][0])
	}

	y.streamList = streams
	return nil
}

func (y *Youtube) getVideoInfo() error {
	url := "http://youtube.com/get_video_info?video_id=" + y.videoID
	log.Printf("url: %s", url)
	resp, err := http.Get(url)
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

func (y *Youtube) findVideoId(url string) error {
	videoId := url
	if strings.Contains(videoId, "youtu") || strings.ContainsAny(videoId, "\"?&/<%=") {
		re_list := []*regexp.Regexp{
			regexp.MustCompile(`(?:v|embed|watch\?v)(?:=|/)([^"&?/=%]{11})`),
			regexp.MustCompile(`(?:=|/)([^"&?/=%]{11})`),
			regexp.MustCompile(`([^"&?/=%]{11})`),
		}
		for _, re := range re_list {
			if is_match := re.MatchString(videoId); is_match {
				subs := re.FindStringSubmatch(videoId)
				videoId = subs[1]
			}
		}
	}
	log.Printf("Found video id: '%s'", videoId)
	y.videoID = videoId
	if strings.ContainsAny(videoId, "?&/<%=") {
		return errors.New("invalid characters in video id")
	}
	if len(videoId) < 10 {
		return errors.New("the video id must be at least 10 characters long")
	}
	return nil
}

func videoDLWorker(destFile string, target string) {
	resp, err := http.Get(target)
	if err != nil {
		log.Printf("Http.Get\nerror: %s\ntarget: %s\n", err, target)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("reading answer: non 200 status code received: '%s'", err)
	}
	out, err := os.Create(destFile)
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Println("download video err=", err)
		return
	}
}
