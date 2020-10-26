package youtube_test

import (
	"io"
	"os"

	"github.com/kkdai/youtube/v2"
)

// ExampleDownload : Example code for how to use this package for download video.
func ExampleClient() {
	videoID := "BaW_jenozKc"
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		panic(err)
	}

	resp, err := client.GetStream(video, &video.Formats[0])
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	file, err := os.Create("video.mp4")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		panic(err)
	}
}
