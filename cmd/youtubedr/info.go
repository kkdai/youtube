package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// Define two new struct in local scope
type VideoFormat struct {
	Itag          int
	FPS           int
	VideoQuality  string
	AudioQuality  string
	AudioChannels int
	Language      string
	Size          int64
	Bitrate       int
	MimeType      string
}

type VideoInfo struct {
	ID          string
	Title       string
	Author      string
	Duration    string
	Description string
	Formats     []VideoFormat
}

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print metadata of the desired video",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkOutputFormat()
	},
	Run: func(cmd *cobra.Command, args []string) {
		video, err := getDownloader().GetVideo(args[0])
		exitOnError(err)

		videoInfo := VideoInfo{
			Title:       video.Title,
			Author:      video.Author,
			Duration:    video.Duration.String(),
			Description: video.Description,
		}

		for _, format := range video.Formats {
			bitrate := format.AverageBitrate
			if bitrate == 0 {
				// Some formats don't have the average bitrate
				bitrate = format.Bitrate
			}

			size := format.ContentLength
			if size == 0 {
				// Some formats don't have this information
				size = int64(float64(bitrate) * video.Duration.Seconds() / 8)
			}

			videoInfo.Formats = append(videoInfo.Formats, VideoFormat{
				Itag:          format.ItagNo,
				FPS:           format.FPS,
				VideoQuality:  format.QualityLabel,
				AudioQuality:  strings.ToLower(strings.TrimPrefix(format.AudioQuality, "AUDIO_QUALITY_")),
				AudioChannels: format.AudioChannels,
				Size:          size,
				Bitrate:       bitrate,
				MimeType:      format.MimeType,
				Language:      format.LanguageDisplayName(),
			})
		}

		exitOnError(writeOutput(os.Stdout, &videoInfo, func(w io.Writer) {
			writeInfoOutput(w, &videoInfo)
		}))
	},
}

func writeInfoOutput(w io.Writer, info *VideoInfo) {
	fmt.Println("Title:      ", info.Title)
	fmt.Println("Author:     ", info.Author)
	fmt.Println("Duration:   ", info.Duration)
	if printDescription {
		fmt.Println("Description:", info.Description)
	}
	fmt.Println()

	table := tablewriter.NewWriter(w)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{
		"itag",
		"fps",
		"video\nquality",
		"audio\nquality",
		"audio\nchannels",
		"size [MB]",
		"bitrate",
		"MimeType",
		"language",
	})

	for _, format := range info.Formats {
		table.Append([]string{
			strconv.Itoa(format.Itag),
			strconv.Itoa(format.FPS),
			format.VideoQuality,
			format.AudioQuality,
			strconv.Itoa(format.AudioChannels),
			fmt.Sprintf("%0.1f", float64(format.Size)/1024/1024),
			strconv.Itoa(format.Bitrate),
			format.MimeType,
			format.Language,
		})
	}

	table.Render()
}

var printDescription bool

func init() {
	rootCmd.AddCommand(infoCmd)
	addFormatFlag(infoCmd.Flags())
	infoCmd.Flags().BoolVarP(&printDescription, "description", "d", false, "Print description")
}
