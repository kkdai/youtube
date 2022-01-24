package youtube

type playerResponseData struct {
	ResponseContext        ResponseContext        `json:"responseContext"`
	TrackingParams         string                 `json:"trackingParams"`
	AdBreakParams          string                 `json:"adBreakParams"`
	PlayabilityStatus      PlayabilityStatus      `json:"playabilityStatus"`
	StreamingData          StreamingData          `json:"streamingData"`
	HeartbeatParams        HeartbeatParams        `json:"heartbeatParams,omitempty"`
	PlaybackTracking       PlaybackTracking       `json:"playbackTracking"`
	VideoDetails           VideoDetails           `json:"videoDetails"`
	Annotations            []Annotations          `json:"annotations"`
	PlayerConfig           PlayerConfig           `json:"playerConfig"`
	Storyboards            Storyboards            `json:"storyboards"`
	Attestation            Attestation            `json:"attestation"`
	PlayerSettingsMenuData PlayerSettingsMenuData `json:"playerSettingsMenuData"`
	Captions               Captions               `json:"captions"`
}
type Params struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type ServiceTrackingParams struct {
	Service string   `json:"service"`
	Params  []Params `json:"params"`
}
type ResponseContext struct {
	VisitorData           string                  `json:"visitorData"`
	MaxAgeSeconds         int                     `json:"maxAgeSeconds"`
	ServiceTrackingParams []ServiceTrackingParams `json:"serviceTrackingParams"`
}
type LiveStreamabilityRenderer struct {
	VideoID     string `json:"videoId"`
	Params      string `json:"params"`
	BroadcastID string `json:"broadcastId"`
	PollDelayMs string `json:"pollDelayMs"`
}
type LiveStreamability struct {
	ButtonRenderer            ButtonRenderer            `json:"buttonRenderer"`
	LiveStreamabilityRenderer LiveStreamabilityRenderer `json:"liveStreamabilityRenderer"`
}
type YpcGetOfflineUpsellEndpoint struct {
	Params string `json:"params"`
}
type ServiceEndpoint struct {
	ClickTrackingParams         string                      `json:"clickTrackingParams"`
	YpcGetOfflineUpsellEndpoint YpcGetOfflineUpsellEndpoint `json:"ypcGetOfflineUpsellEndpoint"`
}
type ButtonRenderer struct {
	ServiceEndpoint ServiceEndpoint `json:"serviceEndpoint"`
	TrackingParams  string          `json:"trackingParams"`
}
type Offlineability struct {
	ButtonRenderer ButtonRenderer `json:"buttonRenderer"`
}
type MiniplayerRenderer struct {
	PlaybackMode string `json:"playbackMode"`
}
type Miniplayer struct {
	MiniplayerRenderer MiniplayerRenderer `json:"miniplayerRenderer"`
}
type PlayabilityStatus struct {
	Status            string            `json:"status"`
	PlayableInEmbed   bool              `json:"playableInEmbed"`
	LiveStreamability LiveStreamability `json:"liveStreamability"`
	Offlineability    Offlineability    `json:"offlineability"`
	Miniplayer        Miniplayer        `json:"miniplayer"`
	ContextParams     string            `json:"contextParams"`
	Reason            string            `json:"reason"`
	Watermark         Watermark         `json:"watermark"`
	ErrorScreen       ErrorScreen       `json:"errorScreen"`
}
type Format struct {
	ItagNo           int    `json:"itag"`
	URL              string `json:"url"`
	Cipher           string `json:"signatureCipher,omitempty"`
	MimeType         string `json:"mimeType"`
	Bitrate          int    `json:"bitrate"`
	Width            int    `json:"width"`
	Height           int    `json:"height"`
	LastModified     string `json:"lastModified"`
	ContentLength    string `json:"contentLength,omitempty"`
	Quality          string `json:"quality"`
	FPS              int    `json:"fps"`
	QualityLabel     string `json:"qualityLabel"`
	ProjectionType   string `json:"projectionType"`
	AverageBitrate   int    `json:"averageBitrate,omitempty"`
	AudioQuality     string `json:"audioQuality"`
	ApproxDurationMs string `json:"approxDurationMs"`
	AudioSampleRate  string `json:"audioSampleRate"`
	AudioChannels    int    `json:"audioChannels"`
}
type InitRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}
type IndexRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}
type ColorInfo struct {
	Primaries               string `json:"primaries"`
	TransferCharacteristics string `json:"transferCharacteristics"`
	MatrixCoefficients      string `json:"matrixCoefficients"`
}
type AdaptiveFormats struct {
	Itag              int        `json:"itag"`
	URL               string     `json:"url"`
	MimeType          string     `json:"mimeType"`
	Bitrate           int        `json:"bitrate"`
	Width             int        `json:"width,omitempty"`
	Height            int        `json:"height,omitempty"`
	InitRange         InitRange  `json:"initRange,omitempty"`
	IndexRange        IndexRange `json:"indexRange,omitempty"`
	LastModified      string     `json:"lastModified"`
	ContentLength     string     `json:"contentLength,omitempty"`
	Quality           string     `json:"quality"`
	FPS               int        `json:"fps,omitempty"`
	QualityLabel      string     `json:"qualityLabel,omitempty"`
	ProjectionType    string     `json:"projectionType"`
	TargetDurationSec int        `json:"targetDurationSec"`
	MaxDvrDurationSec int        `json:"maxDvrDurationSec"`
	AverageBitrate    int        `json:"averageBitrate,omitempty"`
	ApproxDurationMs  string     `json:"approxDurationMs,omitempty"`
	ColorInfo         ColorInfo  `json:"colorInfo,omitempty"`
	HighReplication   bool       `json:"highReplication,omitempty"`
	AudioQuality      string     `json:"audioQuality,omitempty"`
	AudioSampleRate   string     `json:"audioSampleRate,omitempty"`
	AudioChannels     int        `json:"audioChannels,omitempty"`
}
type StreamingData struct {
	ExpiresInSeconds   string            `json:"expiresInSeconds"`
	Formats            []Format          `json:"formats"`
	AdaptiveFormats    []AdaptiveFormats `json:"adaptiveFormats"`
	OnesieStreamingURL string            `json:"onesieStreamingUrl"`
	DashManifestURL    string            `json:"dashManifestUrl"`
	HlsManifestURL     string            `json:"hlsManifestUrl"`
}
type HeartbeatParams struct {
	IntervalMilliseconds string `json:"intervalMilliseconds"`
	SoftFailOnError      bool   `json:"softFailOnError"`
	HeartbeatServerData  string `json:"heartbeatServerData"`
}
type Headers struct {
	HeaderType string `json:"headerType"`
}
type VideostatsPlaybackURL struct {
	BaseURL string    `json:"baseUrl"`
	Headers []Headers `json:"headers"`
}
type VideostatsDelayplayURL struct {
	BaseURL                 string    `json:"baseUrl"`
	ElapsedMediaTimeSeconds int       `json:"elapsedMediaTimeSeconds"`
	Headers                 []Headers `json:"headers"`
}
type VideostatsWatchtimeURL struct {
	BaseURL string    `json:"baseUrl"`
	Headers []Headers `json:"headers"`
}
type PtrackingURL struct {
	BaseURL string    `json:"baseUrl"`
	Headers []Headers `json:"headers"`
}
type QoeURL struct {
	BaseURL string    `json:"baseUrl"`
	Headers []Headers `json:"headers"`
}
type AtrURL struct {
	BaseURL                 string    `json:"baseUrl"`
	ElapsedMediaTimeSeconds int       `json:"elapsedMediaTimeSeconds"`
	Headers                 []Headers `json:"headers"`
}
type EngageURL struct {
	BaseURL string    `json:"baseUrl"`
	Headers []Headers `json:"headers"`
}
type YoutubeRemarketingURL struct {
	BaseURL                 string    `json:"baseUrl"`
	ElapsedMediaTimeSeconds int       `json:"elapsedMediaTimeSeconds"`
	Headers                 []Headers `json:"headers"`
}
type PlaybackTracking struct {
	VideostatsPlaybackURL                   VideostatsPlaybackURL  `json:"videostatsPlaybackUrl"`
	VideostatsDelayplayURL                  VideostatsDelayplayURL `json:"videostatsDelayplayUrl"`
	VideostatsWatchtimeURL                  VideostatsWatchtimeURL `json:"videostatsWatchtimeUrl"`
	PtrackingURL                            PtrackingURL           `json:"ptrackingUrl"`
	QoeURL                                  QoeURL                 `json:"qoeUrl"`
	AtrURL                                  AtrURL                 `json:"atrUrl"`
	EngageURL                               EngageURL              `json:"engageUrl"`
	VideostatsScheduledFlushWalltimeSeconds []int                  `json:"videostatsScheduledFlushWalltimeSeconds"`
	VideostatsDefaultFlushIntervalSeconds   int                    `json:"videostatsDefaultFlushIntervalSeconds"`
	YoutubeRemarketingURL                   YoutubeRemarketingURL  `json:"youtubeRemarketingUrl"`
}
type Thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
type Thumbnails struct {
	Thumbnails []Thumbnail `json:"thumbnails"`
}
type VideoDetails struct {
	VideoID                string     `json:"videoId"`
	Title                  string     `json:"title"`
	LengthSeconds          string     `json:"lengthSeconds"`
	IsLive                 bool       `json:"isLive,omitempty"`
	Keywords               []string   `json:"keywords"`
	ChannelID              string     `json:"channelId"`
	IsOwnerViewing         bool       `json:"isOwnerViewing"`
	ShortDescription       string     `json:"shortDescription"`
	IsCrawlable            bool       `json:"isCrawlable"`
	IsLiveDvrEnabled       bool       `json:"isLiveDvrEnabled,omitempty"`
	Thumbnails             Thumbnails `json:"thumbnail"`
	LiveChunkReadahead     int        `json:"liveChunkReadahead,omitempty"`
	AllowRatings           bool       `json:"allowRatings"`
	ViewCount              string     `json:"viewCount"`
	Author                 string     `json:"author"`
	IsPrivate              bool       `json:"isPrivate"`
	IsUnpluggedCorpus      bool       `json:"isUnpluggedCorpus"`
	IsLiveContent          bool       `json:"isLiveContent"`
	IsLowLatencyLiveStream bool       `json:"isLowLatencyLiveStream,omitempty"`
	LatencyClass           string     `json:"latencyClass,omitempty"`
}
type Watermark struct {
	Thumbnails []Thumbnails `json:"thumbnails"`
}
type FeaturedChannel struct {
	StartTimeMs    string    `json:"startTimeMs"`
	EndTimeMs      string    `json:"endTimeMs"`
	Watermark      Watermark `json:"watermark"`
	TrackingParams string    `json:"trackingParams"`
}
type PlayerAnnotationsExpandedRenderer struct {
	FeaturedChannel   FeaturedChannel `json:"featuredChannel"`
	AllowSwipeDismiss bool            `json:"allowSwipeDismiss"`
}
type Annotations struct {
	PlayerAnnotationsExpandedRenderer PlayerAnnotationsExpandedRenderer `json:"playerAnnotationsExpandedRenderer"`
}
type AudioConfig struct {
	LoudnessDb              float64 `json:"loudnessDb"`
	PerceptualLoudnessDb    float64 `json:"perceptualLoudnessDb"`
	EnablePerFormatLoudness bool    `json:"enablePerFormatLoudness"`
}
type ExoPlayerConfig struct {
	UseExoPlayer                                            bool     `json:"useExoPlayer"`
	UseAdaptiveBitrate                                      bool     `json:"useAdaptiveBitrate"`
	MaxInitialByteRate                                      int      `json:"maxInitialByteRate"`
	MinDurationForQualityIncreaseMs                         int      `json:"minDurationForQualityIncreaseMs"`
	MaxDurationForQualityDecreaseMs                         int      `json:"maxDurationForQualityDecreaseMs"`
	MinDurationToRetainAfterDiscardMs                       int      `json:"minDurationToRetainAfterDiscardMs"`
	LowWatermarkMs                                          int      `json:"lowWatermarkMs"`
	HighWatermarkMs                                         int      `json:"highWatermarkMs"`
	LowPoolLoad                                             float64  `json:"lowPoolLoad"`
	HighPoolLoad                                            float64  `json:"highPoolLoad"`
	SufficientBandwidthOverhead                             float64  `json:"sufficientBandwidthOverhead"`
	BufferChunkSizeKb                                       int      `json:"bufferChunkSizeKb"`
	HTTPConnectTimeoutMs                                    int      `json:"httpConnectTimeoutMs"`
	HTTPReadTimeoutMs                                       int      `json:"httpReadTimeoutMs"`
	NumAudioSegmentsPerFetch                                int      `json:"numAudioSegmentsPerFetch"`
	NumVideoSegmentsPerFetch                                int      `json:"numVideoSegmentsPerFetch"`
	MinDurationForPlaybackStartMs                           int      `json:"minDurationForPlaybackStartMs"`
	EnableExoplayerReuse                                    bool     `json:"enableExoplayerReuse"`
	UseRadioTypeForInitialQualitySelection                  bool     `json:"useRadioTypeForInitialQualitySelection"`
	BlacklistFormatOnError                                  bool     `json:"blacklistFormatOnError"`
	EnableBandaidHTTPDataSource                             bool     `json:"enableBandaidHttpDataSource"`
	HTTPLoadTimeoutMs                                       int      `json:"httpLoadTimeoutMs"`
	CanPlayHdDrm                                            bool     `json:"canPlayHdDrm"`
	VideoBufferSegmentCount                                 int      `json:"videoBufferSegmentCount"`
	AudioBufferSegmentCount                                 int      `json:"audioBufferSegmentCount"`
	UseAbruptSplicing                                       bool     `json:"useAbruptSplicing"`
	MinRetryCount                                           int      `json:"minRetryCount"`
	MinChunksNeededToPreferOffline                          int      `json:"minChunksNeededToPreferOffline"`
	SecondsToMaxAggressiveness                              int      `json:"secondsToMaxAggressiveness"`
	EnableSurfaceviewResizeWorkaround                       bool     `json:"enableSurfaceviewResizeWorkaround"`
	EnableVp9IfThresholdsPass                               bool     `json:"enableVp9IfThresholdsPass"`
	MatchQualityToViewportOnUnfullscreen                    bool     `json:"matchQualityToViewportOnUnfullscreen"`
	LowAudioQualityConnTypes                                []string `json:"lowAudioQualityConnTypes"`
	UseDashForLiveStreams                                   bool     `json:"useDashForLiveStreams"`
	EnableLibvpxVideoTrackRenderer                          bool     `json:"enableLibvpxVideoTrackRenderer"`
	LowAudioQualityBandwidthThresholdBps                    int      `json:"lowAudioQualityBandwidthThresholdBps"`
	EnableVariableSpeedPlayback                             bool     `json:"enableVariableSpeedPlayback"`
	PreferOnesieBufferedFormat                              bool     `json:"preferOnesieBufferedFormat"`
	MinimumBandwidthSampleBytes                             int      `json:"minimumBandwidthSampleBytes"`
	UseDashForOtfAndCompletedLiveStreams                    bool     `json:"useDashForOtfAndCompletedLiveStreams"`
	DisableCacheAwareVideoFormatEvaluation                  bool     `json:"disableCacheAwareVideoFormatEvaluation"`
	UseLiveDvrForDashLiveStreams                            bool     `json:"useLiveDvrForDashLiveStreams"`
	CronetResetTimeoutOnRedirects                           bool     `json:"cronetResetTimeoutOnRedirects"`
	EmitVideoDecoderChangeEvents                            bool     `json:"emitVideoDecoderChangeEvents"`
	OnesieVideoBufferLoadTimeoutMs                          string   `json:"onesieVideoBufferLoadTimeoutMs"`
	OnesieVideoBufferReadTimeoutMs                          string   `json:"onesieVideoBufferReadTimeoutMs"`
	LibvpxEnableGl                                          bool     `json:"libvpxEnableGl"`
	EnableVp9EncryptedIfThresholdsPass                      bool     `json:"enableVp9EncryptedIfThresholdsPass"`
	EnableOpus                                              bool     `json:"enableOpus"`
	UsePredictedBuffer                                      bool     `json:"usePredictedBuffer"`
	MaxReadAheadMediaTimeMs                                 int      `json:"maxReadAheadMediaTimeMs"`
	UseMediaTimeCappedLoadControl                           bool     `json:"useMediaTimeCappedLoadControl"`
	AllowCacheOverrideToLowerQualitiesWithinRange           int      `json:"allowCacheOverrideToLowerQualitiesWithinRange"`
	AllowDroppingUndecodedFrames                            bool     `json:"allowDroppingUndecodedFrames"`
	MinDurationForPlaybackRestartMs                         int      `json:"minDurationForPlaybackRestartMs"`
	ServerProvidedBandwidthHeader                           string   `json:"serverProvidedBandwidthHeader"`
	LiveOnlyPegStrategy                                     string   `json:"liveOnlyPegStrategy"`
	EnableRedirectorHostFallback                            bool     `json:"enableRedirectorHostFallback"`
	EnableHighlyAvailableFormatFallbackOnPcr                bool     `json:"enableHighlyAvailableFormatFallbackOnPcr"`
	RecordTrackRendererTimingEvents                         bool     `json:"recordTrackRendererTimingEvents"`
	MinErrorsForRedirectorHostFallback                      int      `json:"minErrorsForRedirectorHostFallback"`
	NonHardwareMediaCodecNames                              []string `json:"nonHardwareMediaCodecNames"`
	EnableVp9IfInHardware                                   bool     `json:"enableVp9IfInHardware"`
	EnableVp9EncryptedIfInHardware                          bool     `json:"enableVp9EncryptedIfInHardware"`
	UseOpusMedAsLowQualityAudio                             bool     `json:"useOpusMedAsLowQualityAudio"`
	MinErrorsForPcrFallback                                 int      `json:"minErrorsForPcrFallback"`
	UseStickyRedirectHTTPDataSource                         bool     `json:"useStickyRedirectHttpDataSource"`
	OnlyVideoBandwidth                                      bool     `json:"onlyVideoBandwidth"`
	UseRedirectorOnNetworkChange                            bool     `json:"useRedirectorOnNetworkChange"`
	EnableMaxReadaheadAbrThreshold                          bool     `json:"enableMaxReadaheadAbrThreshold"`
	CacheCheckDirectoryWritabilityOnce                      bool     `json:"cacheCheckDirectoryWritabilityOnce"`
	PredictorType                                           string   `json:"predictorType"`
	SlidingPercentile                                       float64  `json:"slidingPercentile"`
	SlidingWindowSize                                       int      `json:"slidingWindowSize"`
	MaxFrameDropIntervalMs                                  int      `json:"maxFrameDropIntervalMs"`
	IgnoreLoadTimeoutForFallback                            bool     `json:"ignoreLoadTimeoutForFallback"`
	ServerBweMultiplier                                     int      `json:"serverBweMultiplier"`
	DrmMaxKeyfetchDelayMs                                   int      `json:"drmMaxKeyfetchDelayMs"`
	MaxResolutionForWhiteNoise                              int      `json:"maxResolutionForWhiteNoise"`
	WhiteNoiseRenderEffectMode                              string   `json:"whiteNoiseRenderEffectMode"`
	EnableLibvpxHdr                                         bool     `json:"enableLibvpxHdr"`
	EnableCacheAwareStreamSelection                         bool     `json:"enableCacheAwareStreamSelection"`
	UseExoCronetDataSource                                  bool     `json:"useExoCronetDataSource"`
	WhiteNoiseScale                                         int      `json:"whiteNoiseScale"`
	WhiteNoiseOffset                                        int      `json:"whiteNoiseOffset"`
	PreventVideoFrameLaggingWithLibvpx                      bool     `json:"preventVideoFrameLaggingWithLibvpx"`
	EnableMediaCodecHdr                                     bool     `json:"enableMediaCodecHdr"`
	EnableMediaCodecSwHdr                                   bool     `json:"enableMediaCodecSwHdr"`
	LiveOnlyWindowChunks                                    int      `json:"liveOnlyWindowChunks"`
	BearerMinDurationToRetainAfterDiscardMs                 []int    `json:"bearerMinDurationToRetainAfterDiscardMs"`
	ForceWidevineL3                                         bool     `json:"forceWidevineL3"`
	UseAverageBitrate                                       bool     `json:"useAverageBitrate"`
	UseMedialibAudioTrackRendererForLive                    bool     `json:"useMedialibAudioTrackRendererForLive"`
	UseExoPlayerV2                                          bool     `json:"useExoPlayerV2"`
	EnableManifestlessResumeVideo                           bool     `json:"enableManifestlessResumeVideo"`
	LogMediaRequestEventsToCsi                              bool     `json:"logMediaRequestEventsToCsi"`
	OnesieFixNonZeroStartTimeFormatSelection                bool     `json:"onesieFixNonZeroStartTimeFormatSelection"`
	LiveOnlyReadaheadStepSizeChunks                         int      `json:"liveOnlyReadaheadStepSizeChunks"`
	LiveOnlyBufferHealthHalfLifeSeconds                     int      `json:"liveOnlyBufferHealthHalfLifeSeconds"`
	LiveOnlyMinBufferHealthRatio                            float64  `json:"liveOnlyMinBufferHealthRatio"`
	LiveOnlyMinLatencyToSeekRatio                           int      `json:"liveOnlyMinLatencyToSeekRatio"`
	ManifestlessPartialChunkStrategy                        string   `json:"manifestlessPartialChunkStrategy"`
	IgnoreViewportSizeWhenSticky                            bool     `json:"ignoreViewportSizeWhenSticky"`
	EnableLibvpxFallback                                    bool     `json:"enableLibvpxFallback"`
	DisableLibvpxLoopFilter                                 bool     `json:"disableLibvpxLoopFilter"`
	EnableVpxMediaView                                      bool     `json:"enableVpxMediaView"`
	HdrMinScreenBrightness                                  int      `json:"hdrMinScreenBrightness"`
	HdrMaxScreenBrightnessThreshold                         int      `json:"hdrMaxScreenBrightnessThreshold"`
	OnesieDataSourceAboveCacheDataSource                    bool     `json:"onesieDataSourceAboveCacheDataSource"`
	HTTPNonplayerLoadTimeoutMs                              int      `json:"httpNonplayerLoadTimeoutMs"`
	NumVideoSegmentsPerFetchStrategy                        string   `json:"numVideoSegmentsPerFetchStrategy"`
	MaxVideoDurationPerFetchMs                              int      `json:"maxVideoDurationPerFetchMs"`
	MaxVideoEstimatedLoadDurationMs                         int      `json:"maxVideoEstimatedLoadDurationMs"`
	EstimatedServerClockHalfLife                            int      `json:"estimatedServerClockHalfLife"`
	EstimatedServerClockStrictOffset                        bool     `json:"estimatedServerClockStrictOffset"`
	MinReadAheadMediaTimeMs                                 int      `json:"minReadAheadMediaTimeMs"`
	ReadAheadGrowthRate                                     int      `json:"readAheadGrowthRate"`
	UseDynamicReadAhead                                     bool     `json:"useDynamicReadAhead"`
	UseYtVodMediaSourceForV2                                bool     `json:"useYtVodMediaSourceForV2"`
	EnableV2Gapless                                         bool     `json:"enableV2Gapless"`
	UseLiveHeadTimeMillis                                   bool     `json:"useLiveHeadTimeMillis"`
	AllowTrackSelectionWithUpdatedVideoItagsForExoV2        bool     `json:"allowTrackSelectionWithUpdatedVideoItagsForExoV2"`
	MaxAllowableTimeBeforeMediaTimeUpdateSec                int      `json:"maxAllowableTimeBeforeMediaTimeUpdateSec"`
	EnableDynamicHdr                                        bool     `json:"enableDynamicHdr"`
	V2PerformEarlyStreamSelection                           bool     `json:"v2PerformEarlyStreamSelection"`
	V2UsePlaybackStreamSelectionResult                      bool     `json:"v2UsePlaybackStreamSelectionResult"`
	V2MinTimeBetweenAbrReevaluationMs                       int      `json:"v2MinTimeBetweenAbrReevaluationMs"`
	AvoidReusePlaybackAcrossLoadvideos                      bool     `json:"avoidReusePlaybackAcrossLoadvideos"`
	EnableInfiniteNetworkLoadingRetries                     bool     `json:"enableInfiniteNetworkLoadingRetries"`
	ReportExoPlayerStateOnTransition                        bool     `json:"reportExoPlayerStateOnTransition"`
	ManifestlessSequenceMethod                              string   `json:"manifestlessSequenceMethod"`
	UseLiveHeadWindow                                       bool     `json:"useLiveHeadWindow"`
	EnableDynamicHdrInHardware                              bool     `json:"enableDynamicHdrInHardware"`
	UltralowAudioQualityBandwidthThresholdBps               int      `json:"ultralowAudioQualityBandwidthThresholdBps"`
	RetryLiveNetNocontentWithDelay                          bool     `json:"retryLiveNetNocontentWithDelay"`
	IgnoreUnneededSeeksToLiveHead                           bool     `json:"ignoreUnneededSeeksToLiveHead"`
	AdaptiveLiveHeadWindow                                  bool     `json:"adaptiveLiveHeadWindow"`
	DrmMetricsQoeLoggingFraction                            float64  `json:"drmMetricsQoeLoggingFraction"`
	LiveNetNocontentMaximumErrors                           int      `json:"liveNetNocontentMaximumErrors"`
	WaitForDrmLicenseBeforeProcessingAndroidStuckBufferfull bool     `json:"waitForDrmLicenseBeforeProcessingAndroidStuckBufferfull"`
	UseTimeSeriesBufferPrediction                           bool     `json:"useTimeSeriesBufferPrediction"`
}
type PlaybackStartConfig struct {
	StartTimeToleranceBeforeMs string `json:"startTimeToleranceBeforeMs"`
}
type AdRequestConfig struct {
	FilterTimeEventsOnDelta               int     `json:"filterTimeEventsOnDelta"`
	UseCriticalExecOnAdsPrep              bool    `json:"useCriticalExecOnAdsPrep"`
	UserCriticalExecOnAdsProcessing       bool    `json:"userCriticalExecOnAdsProcessing"`
	EnableCountdownNextToThumbnailAndroid bool    `json:"enableCountdownNextToThumbnailAndroid"`
	PreskipScalingFactorAndroid           float64 `json:"preskipScalingFactorAndroid"`
	PreskipPaddingAndroid                 int     `json:"preskipPaddingAndroid"`
}
type NetworkProtocolConfig struct {
	UseQuic bool `json:"useQuic"`
}
type AndroidCronetResponsePriority struct {
	PriorityValue string `json:"priorityValue"`
}
type AndroidMetadataNetworkConfig struct {
	CoalesceRequests bool `json:"coalesceRequests"`
}
type AndroidNetworkStackConfig struct {
	NetworkStack                  string                        `json:"networkStack"`
	AndroidCronetResponsePriority AndroidCronetResponsePriority `json:"androidCronetResponsePriority"`
	AndroidMetadataNetworkConfig  AndroidMetadataNetworkConfig  `json:"androidMetadataNetworkConfig"`
}
type LidarSdkConfig struct {
	EnableActiveViewReporter             bool `json:"enableActiveViewReporter"`
	UseMediaTime                         bool `json:"useMediaTime"`
	SendTosMetrics                       bool `json:"sendTosMetrics"`
	UsePlayerState                       bool `json:"usePlayerState"`
	EnableIosAppStateCheck               bool `json:"enableIosAppStateCheck"`
	EnableImprovedSizeReportingAndroid   bool `json:"enableImprovedSizeReportingAndroid"`
	EnableIsAndroidVideoAlwaysMeasurable bool `json:"enableIsAndroidVideoAlwaysMeasurable"`
}
type InitialBandwidthEstimates struct {
	DetailedNetworkType string `json:"detailedNetworkType"`
	BandwidthBps        string `json:"bandwidthBps"`
}
type AndroidMedialibConfig struct {
	IsItag18MainProfile                     bool                        `json:"isItag18MainProfile"`
	DashManifestVersion                     int                         `json:"dashManifestVersion"`
	InitialBandwidthEstimates               []InitialBandwidthEstimates `json:"initialBandwidthEstimates"`
	ViewportSizeFraction                    float64                     `json:"viewportSizeFraction"`
	SelectLowQualityStreamsWithHighBitrates bool                        `json:"selectLowQualityStreamsWithHighBitrates"`
	EnablePrerollPrebuffer                  bool                        `json:"enablePrerollPrebuffer"`
	PrebufferOptimizeForViewportSize        bool                        `json:"prebufferOptimizeForViewportSize"`
	HpqViewportSizeFraction                 float64                     `json:"hpqViewportSizeFraction"`
}
type PlayerControlsConfig struct {
	ShowCachedInTimebar bool `json:"showCachedInTimebar"`
}
type VariableSpeedConfig struct {
	ShowVariableSpeedDisabledDialog bool `json:"showVariableSpeedDisabledDialog"`
}
type DecodeQualityConfig struct {
	MaximumVideoDecodeVerticalResolution int `json:"maximumVideoDecodeVerticalResolution"`
}
type VrConfig struct {
	AllowVr                            bool   `json:"allowVr"`
	AllowSubtitles                     bool   `json:"allowSubtitles"`
	ShowHqButton                       bool   `json:"showHqButton"`
	SphericalDirectionLoggingEnabled   bool   `json:"sphericalDirectionLoggingEnabled"`
	EnableAndroidVr180MagicWindow      bool   `json:"enableAndroidVr180MagicWindow"`
	EnableAndroidMagicWindowEduOverlay bool   `json:"enableAndroidMagicWindowEduOverlay"`
	MagicWindowEduOverlayText          string `json:"magicWindowEduOverlayText"`
	MagicWindowEduOverlayAnimationURL  string `json:"magicWindowEduOverlayAnimationUrl"`
	EnableMagicWindowZoom              bool   `json:"enableMagicWindowZoom"`
}
type QoeStatsClientConfig struct {
	BatchedEntriesPeriodMs string `json:"batchedEntriesPeriodMs"`
}
type AndroidPlayerStatsConfig struct {
	UsePblForAttestationReporting      bool `json:"usePblForAttestationReporting"`
	UsePblForHeartbeatReporting        bool `json:"usePblForHeartbeatReporting"`
	UsePblForPlaybacktrackingReporting bool `json:"usePblForPlaybacktrackingReporting"`
	UsePblForQoeReporting              bool `json:"usePblForQoeReporting"`
	ChangeCpnOnFatalPlaybackError      bool `json:"changeCpnOnFatalPlaybackError"`
}
type StickyQualitySelectionConfig struct {
	StickySelectionType                                  string `json:"stickySelectionType"`
	ExpirationTimeSinceLastManualVideoQualitySelectionMs string `json:"expirationTimeSinceLastManualVideoQualitySelectionMs"`
	ExpirationTimeSinceLastPlaybackStartMs               string `json:"expirationTimeSinceLastPlaybackStartMs"`
	StickyCeilingOverridesSimpleBitrateCap               bool   `json:"stickyCeilingOverridesSimpleBitrateCap"`
}
type AdSurveyRequestConfig struct {
	UseGetRequests bool `json:"useGetRequests"`
}
type LivePlayerConfig struct {
	LiveReadaheadSeconds  int `json:"liveReadaheadSeconds"`
	LiveHeadWindowSeconds int `json:"liveHeadWindowSeconds"`
}
type RetryConfig struct {
	RetryEligibleErrors                     []string `json:"retryEligibleErrors"`
	RetryUnderSameConditionAttempts         int      `json:"retryUnderSameConditionAttempts"`
	RetryWithNewSurfaceAttempts             int      `json:"retryWithNewSurfaceAttempts"`
	ProgressiveFallbackOnNonNetworkErrors   bool     `json:"progressiveFallbackOnNonNetworkErrors"`
	L3FallbackOnDrmErrors                   bool     `json:"l3FallbackOnDrmErrors"`
	RetryAfterCacheRemoval                  bool     `json:"retryAfterCacheRemoval"`
	WidevineL3EnforcedFallbackOnDrmErrors   bool     `json:"widevineL3EnforcedFallbackOnDrmErrors"`
	ExoProxyableFormatFallback              bool     `json:"exoProxyableFormatFallback"`
	MaxPlayerRetriesWhenNetworkUnavailable  int      `json:"maxPlayerRetriesWhenNetworkUnavailable"`
	RetryWithLibvpx                         bool     `json:"retryWithLibvpx"`
	SuppressFatalErrorAfterStop             bool     `json:"suppressFatalErrorAfterStop"`
	FallbackFromHfrToSfrOnFormatDecodeError bool     `json:"fallbackFromHfrToSfrOnFormatDecodeError"`
}
type CmsPathProbeConfig struct {
	CmsPathProbeDelayMs int `json:"cmsPathProbeDelayMs"`
}
type DynamicReadaheadConfig struct {
	MaxReadAheadMediaTimeMs             int  `json:"maxReadAheadMediaTimeMs"`
	MinReadAheadMediaTimeMs             int  `json:"minReadAheadMediaTimeMs"`
	ReadAheadGrowthRateMs               int  `json:"readAheadGrowthRateMs"`
	ReadAheadWatermarkMarginRatio       int  `json:"readAheadWatermarkMarginRatio"`
	MinReadAheadWatermarkMarginMs       int  `json:"minReadAheadWatermarkMarginMs"`
	MaxReadAheadWatermarkMarginMs       int  `json:"maxReadAheadWatermarkMarginMs"`
	ShouldIncorporateNetworkActiveState bool `json:"shouldIncorporateNetworkActiveState"`
}
type MediaUstreamerRequestConfig struct {
	EnableVideoPlaybackRequest       bool   `json:"enableVideoPlaybackRequest"`
	VideoPlaybackUstreamerConfig     string `json:"videoPlaybackUstreamerConfig"`
	VideoPlaybackPostEmptyBody       bool   `json:"videoPlaybackPostEmptyBody"`
	IsVideoPlaybackRequestIdempotent bool   `json:"isVideoPlaybackRequestIdempotent"`
}
type PredictedReadaheadConfig struct {
	MinReadaheadMs int `json:"minReadaheadMs"`
	MaxReadaheadMs int `json:"maxReadaheadMs"`
}
type MediaFetchRetryConfig struct {
	InitialDelayMs int     `json:"initialDelayMs"`
	BackoffFactor  float64 `json:"backoffFactor"`
	MaximumDelayMs int     `json:"maximumDelayMs"`
	JitterFactor   float64 `json:"jitterFactor"`
}
type NextRequestPolicy struct {
	TargetAudioReadaheadMs int `json:"targetAudioReadaheadMs"`
	TargetVideoReadaheadMs int `json:"targetVideoReadaheadMs"`
}
type ServerReadaheadConfig struct {
	Enable            bool              `json:"enable"`
	NextRequestPolicy NextRequestPolicy `json:"nextRequestPolicy"`
}
type MediaCommonConfig struct {
	DynamicReadaheadConfig         DynamicReadaheadConfig      `json:"dynamicReadaheadConfig"`
	MediaUstreamerRequestConfig    MediaUstreamerRequestConfig `json:"mediaUstreamerRequestConfig"`
	PredictedReadaheadConfig       PredictedReadaheadConfig    `json:"predictedReadaheadConfig"`
	MediaFetchRetryConfig          MediaFetchRetryConfig       `json:"mediaFetchRetryConfig"`
	MediaFetchMaximumServerErrors  int                         `json:"mediaFetchMaximumServerErrors"`
	MediaFetchMaximumNetworkErrors int                         `json:"mediaFetchMaximumNetworkErrors"`
	MediaFetchMaximumErrors        int                         `json:"mediaFetchMaximumErrors"`
	ServerReadaheadConfig          ServerReadaheadConfig       `json:"serverReadaheadConfig"`
}
type PlayerGestureConfig struct {
	DownAndOutLandscapeAllowed bool `json:"downAndOutLandscapeAllowed"`
	DownAndOutPortraitAllowed  bool `json:"downAndOutPortraitAllowed"`
}
type PlayerConfig struct {
	AudioConfig                  AudioConfig                  `json:"audioConfig"`
	ExoPlayerConfig              ExoPlayerConfig              `json:"exoPlayerConfig"`
	PlaybackStartConfig          PlaybackStartConfig          `json:"playbackStartConfig"`
	AdRequestConfig              AdRequestConfig              `json:"adRequestConfig"`
	NetworkProtocolConfig        NetworkProtocolConfig        `json:"networkProtocolConfig"`
	AndroidNetworkStackConfig    AndroidNetworkStackConfig    `json:"androidNetworkStackConfig"`
	LidarSdkConfig               LidarSdkConfig               `json:"lidarSdkConfig"`
	AndroidMedialibConfig        AndroidMedialibConfig        `json:"androidMedialibConfig"`
	PlayerControlsConfig         PlayerControlsConfig         `json:"playerControlsConfig"`
	VariableSpeedConfig          VariableSpeedConfig          `json:"variableSpeedConfig"`
	DecodeQualityConfig          DecodeQualityConfig          `json:"decodeQualityConfig"`
	VrConfig                     VrConfig                     `json:"vrConfig"`
	QoeStatsClientConfig         QoeStatsClientConfig         `json:"qoeStatsClientConfig"`
	AndroidPlayerStatsConfig     AndroidPlayerStatsConfig     `json:"androidPlayerStatsConfig"`
	StickyQualitySelectionConfig StickyQualitySelectionConfig `json:"stickyQualitySelectionConfig"`
	AdSurveyRequestConfig        AdSurveyRequestConfig        `json:"adSurveyRequestConfig"`
	LivePlayerConfig             LivePlayerConfig             `json:"livePlayerConfig,omitempty"`
	RetryConfig                  RetryConfig                  `json:"retryConfig"`
	CmsPathProbeConfig           CmsPathProbeConfig           `json:"cmsPathProbeConfig"`
	MediaCommonConfig            MediaCommonConfig            `json:"mediaCommonConfig"`
	PlayerGestureConfig          PlayerGestureConfig          `json:"playerGestureConfig"`
}
type PlayerStoryboardSpecRenderer struct {
	Spec             string `json:"spec"`
	RecommendedLevel int    `json:"recommendedLevel"`
}
type Storyboards struct {
	PlayerStoryboardSpecRenderer PlayerStoryboardSpecRenderer `json:"playerStoryboardSpecRenderer"`
}
type PlayerAttestationRenderer struct {
	Challenge string `json:"challenge"`
}
type Attestation struct {
	PlayerAttestationRenderer PlayerAttestationRenderer `json:"playerAttestationRenderer"`
}
type Visibility struct {
	Types string `json:"types"`
}
type LoggingDirectives struct {
	TrackingParams string     `json:"trackingParams"`
	Visibility     Visibility `json:"visibility"`
}
type PlayerSettingsMenuData struct {
	LoggingDirectives LoggingDirectives `json:"loggingDirectives"`
}

type Runs struct {
	Text string `json:"text"`
}

type Name struct {
	Runs []Runs `json:"runs"`
}

type CaptionTracks struct {
	BaseURL        string `json:"baseUrl"`
	Name           Name   `json:"name"`
	VssID          string `json:"vssId"`
	LanguageCode   string `json:"languageCode"`
	Kind           string `json:"kind"`
	IsTranslatable bool   `json:"isTranslatable"`
}

type AudioTracks struct {
	CaptionTrackIndices []int `json:"captionTrackIndices"`
}

type PlayerCaptionsTracklistRenderer struct {
	CaptionTracks          []CaptionTracks `json:"captionTracks"`
	AudioTracks            []AudioTracks   `json:"audioTracks"`
	DefaultAudioTrackIndex int             `json:"defaultAudioTrackIndex"`
}

type Captions struct {
	PlayerCaptionsTracklistRenderer PlayerCaptionsTracklistRenderer `json:"playerCaptionsTracklistRenderer"`
}

type PlayerErrorMessageRenderer struct {
	Reason    Reason    `json:"reason"`
	Thumbnail Thumbnail `json:"thumbnail"`
}
type ErrorScreen struct {
	PlayerErrorMessageRenderer PlayerErrorMessageRenderer `json:"playerErrorMessageRenderer"`
}

type Reason struct {
	Runs []Runs `json:"runs"`
}
