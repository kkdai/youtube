package youtube_test

import (
	"context"
	"io"
	"os"

	"github.com/kkdai/youtube"
)

//ExampleDownload : Example code for how to use this package for download video.
func ExampleClient() {
	videoID := "BaW_jenozKc"
	ctx := context.Background()
	client := youtube.Client{}

	video, err := client.GetVideoContext(ctx, videoID)
	if err != nil {
		panic(err)
	}

	resp, err := client.Download(ctx, video, &video.Streams[0])
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
