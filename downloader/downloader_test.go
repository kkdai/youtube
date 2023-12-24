package downloader

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kkdai/youtube/v2"
)

var testDownloader = func() (dl Downloader) {
	dl.OutputDir = "download_test"

	return
}()

func TestMain(m *testing.M) {
	exitCode := m.Run()
	// the following code doesn't work under debugger, please delete download files manually
	if err := os.RemoveAll(testDownloader.OutputDir); err != nil {
		panic(err)
	}
	os.Exit(exitCode)
}

func TestDownload_FirstStream(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	ctx := context.Background()

	// youtube-dl test video
	video, err := testDownloader.Client.GetVideoContext(ctx, "BaW_jenozKc")
	require.NoError(err)
	require.NotNil(video)

	assert.Equal(`youtube-dl test video "'/\ä↭𝕐`, video.Title)
	assert.Equal(`Philipp Hagemeister`, video.Author)
	assert.Equal(10*time.Second, video.Duration)
	assert.GreaterOrEqual(len(video.Formats), 18)

	if assert.Greater(len(video.Formats), 0) {
		assert.NoError(testDownloader.Download(ctx, video, &video.Formats[0], ""))
	}
}

func TestYoutube_DownloadWithHighQualityFails(t *testing.T) {
	tests := []struct {
		name    string
		formats []youtube.Format
		message string
	}{
		{
			name:    "video format not found",
			formats: []youtube.Format{{ItagNo: 140}},
			message: "no video format found after filtering",
		},
		{
			name:    "audio format not found",
			formats: []youtube.Format{{ItagNo: 137, Quality: "hd1080", MimeType: "video/mp4", AudioChannels: 0}},
			message: "no audio format found after filtering",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			video := &youtube.Video{
				Formats: tt.formats,
			}

			err := testDownloader.DownloadComposite(context.Background(), "", video, "hd1080", "", "")
			assert.EqualError(t, err, tt.message)
		})
	}
}

func Test_getVideoAudioFormats(t *testing.T) {
	require := require.New(t)

	v := &youtube.Video{Formats: []youtube.Format{
		{ItagNo: 22, MimeType: "video/mp4; codecs=\"avc1.64001F, mp4a.40.2\"", Quality: "hd720", Bitrate: 644187, FPS: 30, Width: 1280, Height: 720, LastModified: "1608827694876411", ContentLength: 0, QualityLabel: "720p", ProjectionType: "RECTANGULAR", AverageBitrate: 0, AudioQuality: "AUDIO_QUALITY_MEDIUM", ApproxDurationMs: "3553976", AudioSampleRate: "44100", AudioChannels: 2},
		{ItagNo: 136, MimeType: "video/mp4; codecs=\"avc1.4d401f\"", Quality: "hd720", Bitrate: 1640619, FPS: 30, Width: 1280, Height: 720, LastModified: "1608827567796013", ContentLength: 228884308, QualityLabel: "720p", ProjectionType: "RECTANGULAR", AverageBitrate: 515229, AudioQuality: "", ApproxDurationMs: "3553899", AudioSampleRate: "", AudioChannels: 0},
		{ItagNo: 398, MimeType: "video/mp4; codecs=\"av01.0.05M.08\"", Quality: "hd720", Bitrate: 1480681, FPS: 30, Width: 1280, Height: 720, LastModified: "1610817401605660", ContentLength: 284974685, QualityLabel: "720p", ProjectionType: "RECTANGULAR", AverageBitrate: 641491, AudioQuality: "", ApproxDurationMs: "3553899", AudioSampleRate: "", AudioChannels: 0},
		{ItagNo: 135, MimeType: "video/mp4; codecs=\"avc1.4d401f\"", Quality: "large", Bitrate: 883476, FPS: 30, Width: 854, Height: 480, LastModified: "1608827567791454", ContentLength: 154191496, QualityLabel: "480p", ProjectionType: "RECTANGULAR", AverageBitrate: 347092, AudioQuality: "", ApproxDurationMs: "3553899", AudioSampleRate: "", AudioChannels: 0},
		{ItagNo: 244, MimeType: "video/webm; codecs=\"vp9\"", Quality: "large", Bitrate: 758804, FPS: 30, Width: 854, Height: 480, LastModified: "1540471378691303", ContentLength: 204421502, QualityLabel: "480p", ProjectionType: "RECTANGULAR", AverageBitrate: 460162, AudioQuality: "", ApproxDurationMs: "3553899", AudioSampleRate: "", AudioChannels: 0},
		{ItagNo: 397, MimeType: "video/mp4; codecs=\"av01.0.04M.08\"", Quality: "large", Bitrate: 689144, FPS: 30, Width: 854, Height: 480, LastModified: "1610818826299370", ContentLength: 156644042, QualityLabel: "480p", ProjectionType: "RECTANGULAR", AverageBitrate: 352613, AudioQuality: "", ApproxDurationMs: "3553899", AudioSampleRate: "", AudioChannels: 0},
		{ItagNo: 18, MimeType: "video/mp4; codecs=\"avc1.42001E, mp4a.40.2\"", Quality: "medium", Bitrate: 585374, FPS: 30, Width: 640, Height: 360, LastModified: "1540458716874187", ContentLength: 260045146, QualityLabel: "360p", ProjectionType: "RECTANGULAR", AverageBitrate: 585361, AudioQuality: "AUDIO_QUALITY_LOW", ApproxDurationMs: "3553976", AudioSampleRate: "44100", AudioChannels: 2},
		{ItagNo: 134, MimeType: "video/mp4; codecs=\"avc1.4d401e\"", Quality: "medium", Bitrate: 562611, FPS: 30, Width: 640, Height: 360, LastModified: "1608827567788287", ContentLength: 107193162, QualityLabel: "360p", ProjectionType: "RECTANGULAR", AverageBitrate: 241297, AudioQuality: "", ApproxDurationMs: "3553899", AudioSampleRate: "", AudioChannels: 0},
		{ItagNo: 243, MimeType: "video/webm; codecs=\"vp9\"", Quality: "medium", Bitrate: 415910, FPS: 30, Width: 640, Height: 360, LastModified: "1540470494540739", ContentLength: 124725362, QualityLabel: "360p", ProjectionType: "RECTANGULAR", AverageBitrate: 280762, AudioQuality: "", ApproxDurationMs: "3553899", AudioSampleRate: "", AudioChannels: 0},
		{ItagNo: 396, MimeType: "video/mp4; codecs=\"av01.0.01M.08\"", Quality: "medium", Bitrate: 360497, FPS: 30, Width: 640, Height: 360, LastModified: "1610815679579168", ContentLength: 97717429, QualityLabel: "360p", ProjectionType: "RECTANGULAR", AverageBitrate: 219966, AudioQuality: "", ApproxDurationMs: "3553899", AudioSampleRate: "", AudioChannels: 0},
		{ItagNo: 133, MimeType: "video/mp4; codecs=\"avc1.4d4015\"", Quality: "small", Bitrate: 275124, FPS: 30, Width: 426, Height: 240, LastModified: "1608827567783949", ContentLength: 58518802, QualityLabel: "240p", ProjectionType: "RECTANGULAR", AverageBitrate: 131728, AudioQuality: "", ApproxDurationMs: "3553899", AudioSampleRate: "", AudioChannels: 0},
		{ItagNo: 242, MimeType: "video/webm; codecs=\"vp9\"", Quality: "small", Bitrate: 226276, FPS: 30, Width: 426, Height: 240, LastModified: "1540470499854365", ContentLength: 71003227, QualityLabel: "240p", ProjectionType: "RECTANGULAR", AverageBitrate: 159831, AudioQuality: "", ApproxDurationMs: "3553899", AudioSampleRate: "", AudioChannels: 0},
		{ItagNo: 395, MimeType: "video/mp4; codecs=\"av01.0.00M.08\"", Quality: "small", Bitrate: 174542, FPS: 30, Width: 426, Height: 240, LastModified: "1610814314861797", ContentLength: 53665816, QualityLabel: "240p", ProjectionType: "RECTANGULAR", AverageBitrate: 120804, AudioQuality: "", ApproxDurationMs: "3553899", AudioSampleRate: "", AudioChannels: 0},
		{ItagNo: 160, MimeType: "video/mp4; codecs=\"avc1.4d400c\"", Quality: "tiny", Bitrate: 112782, FPS: 30, Width: 256, Height: 144, LastModified: "1608827567783792", ContentLength: 34559820, QualityLabel: "144p", ProjectionType: "RECTANGULAR", AverageBitrate: 77795, AudioQuality: "", ApproxDurationMs: "3553899", AudioSampleRate: "", AudioChannels: 0},
		{ItagNo: 278, MimeType: "video/webm; codecs=\"vp9\"", Quality: "tiny", Bitrate: 110159, FPS: 30, Width: 256, Height: 144, LastModified: "1540470389548275", ContentLength: 40870406, QualityLabel: "144p", ProjectionType: "RECTANGULAR", AverageBitrate: 92001, AudioQuality: "", ApproxDurationMs: "3553899", AudioSampleRate: "", AudioChannels: 0},
		{ItagNo: 394, MimeType: "video/mp4; codecs=\"av01.0.00M.08\"", Quality: "tiny", Bitrate: 87070, FPS: 30, Width: 256, Height: 144, LastModified: "1610813435452233", ContentLength: 29475328, QualityLabel: "144p", ProjectionType: "RECTANGULAR", AverageBitrate: 66350, AudioQuality: "", ApproxDurationMs: "3553899", AudioSampleRate: "", AudioChannels: 0},
		{ItagNo: 140, MimeType: "audio/mp4; codecs=\"mp4a.40.2\"", Quality: "tiny", Bitrate: 133909, FPS: 0, Width: 0, Height: 0, LastModified: "1608826185990888", ContentLength: 57518278, QualityLabel: "", ProjectionType: "RECTANGULAR", AverageBitrate: 129473, AudioQuality: "AUDIO_QUALITY_MEDIUM", ApproxDurationMs: "3553976", AudioSampleRate: "44100", AudioChannels: 2},
		{ItagNo: 251, MimeType: "audio/webm; codecs=\"opus\"", Quality: "tiny", Bitrate: 168872, FPS: 0, Width: 0, Height: 0, LastModified: "1540474785125205", ContentLength: 62202897, QualityLabel: "", ProjectionType: "RECTANGULAR", AverageBitrate: 140020, AudioQuality: "AUDIO_QUALITY_MEDIUM", ApproxDurationMs: "3553941", AudioSampleRate: "48000", AudioChannels: 2},
		{ItagNo: 250, MimeType: "audio/webm; codecs=\"opus\"", Quality: "tiny", Bitrate: 92040, FPS: 0, Width: 0, Height: 0, LastModified: "1540474783287362", ContentLength: 32290649, QualityLabel: "", ProjectionType: "RECTANGULAR", AverageBitrate: 72686, AudioQuality: "AUDIO_QUALITY_LOW", ApproxDurationMs: "3553941", AudioSampleRate: "48000", AudioChannels: 2},
		{ItagNo: 249, MimeType: "audio/webm; codecs=\"opus\"", Quality: "tiny", Bitrate: 72862, FPS: 0, Width: 0, Height: 0, LastModified: "1540474783513282", ContentLength: 24839529, QualityLabel: "", ProjectionType: "RECTANGULAR", AverageBitrate: 55914, AudioQuality: "AUDIO_QUALITY_LOW", ApproxDurationMs: "3553941", AudioSampleRate: "48000", AudioChannels: 2},
	}}
	{
		videoFormat, audioFormat, err := getVideoAudioFormats(v, "hd720", "mp4", "")
		require.NoError(err)
		require.NotNil(videoFormat)
		require.Equal(398, videoFormat.ItagNo)
		require.NotNil(audioFormat)
		require.Equal(140, audioFormat.ItagNo)
	}

	{
		videoFormat, audioFormat, err := getVideoAudioFormats(v, "large", "webm", "")
		require.NoError(err)
		require.NotNil(videoFormat)
		require.Equal(244, videoFormat.ItagNo)
		require.NotNil(audioFormat)
		require.Equal(251, audioFormat.ItagNo)
	}
}
