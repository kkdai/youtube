package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// infoJsonCmd represents the info command but output it as JSON for other application to read
var infoJSONCmd = &cobra.Command{
	Use:   "infojson",
	Short: "Print metadata of the desired video in json format",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//Define two new struct in local scope
		type VideoFormat struct {
			Itag         int
			VideoQuality string
			AudioQuality string
			Size         float64
			Bitrate      int
			MimeType     string
		}

		type VideoInfo struct {
			Title        string
			Author       string
			Duration     string
			Description  string
			VideoFormats []VideoFormat
		}
		video, err := getDownloader().GetVideo(args[0])
		exitOnError(err)
		thisVideoFormats := []VideoFormat{}
		for _, format := range video.Formats {

			bitrate := format.AverageBitrate
			if bitrate == 0 {
				// Some formats don't have the average bitrate
				bitrate = format.Bitrate
			}

			size, _ := strconv.ParseInt(format.ContentLength, 10, 64)
			if size == 0 {
				// Some formats don't have this information
				size = int64(float64(bitrate) * video.Duration.Seconds() / 8)
			}

			thisVideoFormats = append(thisVideoFormats, VideoFormat{
				Itag:         format.ItagNo,
				VideoQuality: format.QualityLabel,
				AudioQuality: strings.ToLower(strings.TrimPrefix(format.AudioQuality, "AUDIO_QUALITY_")),
				Size:         float64(size) / 1024 / 1024,
				Bitrate:      bitrate,
				MimeType:     format.MimeType,
			})
		}

		//Prase the output struct
		videoInfo := VideoInfo{
			Title:        video.Title,
			Author:       video.Author,
			Duration:     video.Duration.String(),
			Description:  video.Description,
			VideoFormats: thisVideoFormats,
		}

		//Output it as json
		js, _ := json.MarshalIndent(videoInfo, "", "  ")
		fmt.Println(string(js))
	},
}

func init() {
	rootCmd.AddCommand(infoJSONCmd)
}
