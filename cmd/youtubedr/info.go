package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/kkdai/youtube/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var infoCmdOpts struct {
	outputFormat string
}

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:          "info",
	Short:        "Print metadata of the desired videos",
	Example:      `info https://www.youtube.com/watch\?v\=XbNghLqsVwU https://www.youtube.com/watch\?v\=XbNghLqsVwU`,
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !strings.Contains("|full|media|media-csv|", fmt.Sprintf("|%s|", infoCmdOpts.outputFormat)) {
			return fmt.Errorf("output format %s is not valid", infoCmdOpts.outputFormat)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		for _, videoURL := range args {
			video, err := getDownloader().GetVideo(videoURL)
			if err != nil {
				if infoCmdOpts.outputFormat == "media-csv" {
					fmt.Printf("-2,%s,%s\n", videoURL, err)
				} else {
					fmt.Printf("ERROR: %s - %s\n", videoURL, err)
				}
				continue
			}
			if len(video.Formats) == 0 {
				if infoCmdOpts.outputFormat == "media-csv" {
					fmt.Printf("-3,%s,%s\n", videoURL, "no formats found")
				} else {
					fmt.Printf("ERROR: %s - %s\n", videoURL, "no formats found")
				}
				continue
			} else if len(codec) > 0 {
				filterCodecs(video, codec)
			}

			data := buildFormats(video)

			if infoCmdOpts.outputFormat == "media-csv" {
				fmt.Printf("0,%s,%s\n", videoURL, video.Title)
				w := csv.NewWriter(os.Stdout)
				for _, record := range data {
					if err := w.Write(record); err != nil {
						fmt.Printf("-1,%s,%s\n", videoURL, err)
					}
				}
				w.Flush()
				if err := w.Error(); err != nil {
					fmt.Printf("-1,%s,%s\n", videoURL, err)
					continue
				}
			}

			if infoCmdOpts.outputFormat == "full" || infoCmdOpts.outputFormat == "media" {
				fmt.Println("Title:      ", video.Title)
				if infoCmdOpts.outputFormat == "full" {
					fmt.Println("Author:     ", video.Author)
					fmt.Println("Duration:   ", video.Duration)
					fmt.Println("Description:", video.Description)
					fmt.Println()
				}
				table := tablewriter.NewWriter(os.Stdout)
				table.SetAutoWrapText(false)
				table.SetHeader([]string{"itag", "video quality", "audio quality", "size [MB]", "bitrate", "MimeType"})
				table.AppendBulk(data)
				table.Render()
			}
		}
	},
}

// filterCodecs filters out codec with strings.Contains and uses AND operator
func filterCodecs(video *youtube.Video, codec []string) {
	var formats youtube.FormatList
VideoFormat:
	for _, f := range video.Formats {
		for _, c := range codec {
			if !strings.Contains(f.MimeType, c) {
				continue VideoFormat
			}
		}
		formats = append(formats, f)
	}
	video.Formats = formats
	sort.SliceStable(video.Formats, video.SortBitrateDesc)
}

func buildFormats(video *youtube.Video) (data [][]string) {
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

		data = append(data, []string{
			strconv.Itoa(format.ItagNo),
			format.QualityLabel,
			strings.ToLower(strings.TrimPrefix(format.AudioQuality, "AUDIO_QUALITY_")),
			fmt.Sprintf("%0.1f", float64(size)/1024/1024),
			strconv.Itoa(bitrate),
			format.MimeType,
		})
	}

	return
}

func init() {
	rootCmd.AddCommand(infoCmd)

	infoCmd.Flags().StringVarP(&infoCmdOpts.outputFormat, "output", "o", "full", "full, media, media-csv")
	addCodecFlag(infoCmd.Flags())
}
