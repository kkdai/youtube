package youtube

type PlayerResponseData struct {
	PlayabilityStatus struct {
		Status          string `json:"status"`
		Reason          string `json:"reason"`
		PlayableInEmbed bool   `json:"playableInEmbed"`
		ContextParams   string `json:"contextParams"`
	} `json:"playabilityStatus"`
	StreamingData struct {
		ExpiresInSeconds string `json:"expiresInSeconds"`
		Formats          []struct {
			Stream
			Bitrate          int    `json:"bitrate"`
			Width            int    `json:"width"`
			Height           int    `json:"height"`
			LastModified     string `json:"lastModified"`
			ContentLength    string `json:"contentLength,omitempty"`
			QualityLabel     string `json:"qualityLabel"`
			ProjectionType   string `json:"projectionType"`
			AverageBitrate   int    `json:"averageBitrate,omitempty"`
			AudioQuality     string `json:"audioQuality"`
			ApproxDurationMs string `json:"approxDurationMs"`
			AudioSampleRate  string `json:"audioSampleRate"`
			AudioChannels    int    `json:"audioChannels"`
		} `json:"formats"`
		AdaptiveFormats []struct {
			Stream
			Bitrate   int `json:"bitrate"`
			Width     int `json:"width,omitempty"`
			Height    int `json:"height,omitempty"`
			InitRange struct {
				Start string `json:"start"`
				End   string `json:"end"`
			} `json:"initRange"`
			IndexRange struct {
				Start string `json:"start"`
				End   string `json:"end"`
			} `json:"indexRange"`
			LastModified     string `json:"lastModified"`
			ContentLength    string `json:"contentLength"`
			Fps              int    `json:"fps,omitempty"`
			QualityLabel     string `json:"qualityLabel,omitempty"`
			ProjectionType   string `json:"projectionType"`
			AverageBitrate   int    `json:"averageBitrate"`
			ApproxDurationMs string `json:"approxDurationMs"`
			ColorInfo        struct {
				Primaries               string `json:"primaries"`
				TransferCharacteristics string `json:"transferCharacteristics"`
				MatrixCoefficients      string `json:"matrixCoefficients"`
			} `json:"colorInfo,omitempty"`
			HighReplication bool   `json:"highReplication,omitempty"`
			AudioQuality    string `json:"audioQuality,omitempty"`
			AudioSampleRate string `json:"audioSampleRate,omitempty"`
			AudioChannels   int    `json:"audioChannels,omitempty"`
		} `json:"adaptiveFormats"`
	} `json:"streamingData"`
	Captions struct {
		PlayerCaptionsRenderer struct {
			BaseURL    string `json:"baseUrl"`
			Visibility string `json:"visibility"`
		} `json:"playerCaptionsRenderer"`
		PlayerCaptionsTracklistRenderer struct {
			CaptionTracks []struct {
				BaseURL string `json:"baseUrl"`
				Name    struct {
					SimpleText string `json:"simpleText"`
				} `json:"name"`
				VssID          string `json:"vssId"`
				LanguageCode   string `json:"languageCode"`
				Kind           string `json:"kind"`
				IsTranslatable bool   `json:"isTranslatable"`
			} `json:"captionTracks"`
			AudioTracks []struct {
				CaptionTrackIndices []int `json:"captionTrackIndices"`
			} `json:"audioTracks"`
			TranslationLanguages []struct {
				LanguageCode string `json:"languageCode"`
				LanguageName struct {
					SimpleText string `json:"simpleText"`
				} `json:"languageName"`
			} `json:"translationLanguages"`
			DefaultAudioTrackIndex int `json:"defaultAudioTrackIndex"`
		} `json:"playerCaptionsTracklistRenderer"`
	} `json:"captions"`
	VideoDetails struct {
		VideoID          string `json:"videoId"`
		Title            string `json:"title"`
		LengthSeconds    string `json:"lengthSeconds"`
		ChannelID        string `json:"channelId"`
		IsOwnerViewing   bool   `json:"isOwnerViewing"`
		ShortDescription string `json:"shortDescription"`
		IsCrawlable      bool   `json:"isCrawlable"`
		Thumbnail        struct {
			Thumbnails []struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"thumbnails"`
		} `json:"thumbnail"`
		AverageRating     float64 `json:"averageRating"`
		AllowRatings      bool    `json:"allowRatings"`
		ViewCount         string  `json:"viewCount"`
		Author            string  `json:"author"`
		IsPrivate         bool    `json:"isPrivate"`
		IsUnpluggedCorpus bool    `json:"isUnpluggedCorpus"`
		IsLiveContent     bool    `json:"isLiveContent"`
	} `json:"videoDetails"`
	PlayerConfig struct {
		AudioConfig struct {
			LoudnessDb              float64 `json:"loudnessDb"`
			PerceptualLoudnessDb    float64 `json:"perceptualLoudnessDb"`
			EnablePerFormatLoudness bool    `json:"enablePerFormatLoudness"`
		} `json:"audioConfig"`
		StreamSelectionConfig struct {
			MaxBitrate string `json:"maxBitrate"`
		} `json:"streamSelectionConfig"`
		MediaCommonConfig struct {
			DynamicReadaheadConfig struct {
				MaxReadAheadMediaTimeMs int `json:"maxReadAheadMediaTimeMs"`
				MinReadAheadMediaTimeMs int `json:"minReadAheadMediaTimeMs"`
				ReadAheadGrowthRateMs   int `json:"readAheadGrowthRateMs"`
			} `json:"dynamicReadaheadConfig"`
		} `json:"mediaCommonConfig"`
	} `json:"playerConfig"`
	Storyboards struct {
		PlayerStoryboardSpecRenderer struct {
			Spec string `json:"spec"`
		} `json:"playerStoryboardSpecRenderer"`
	} `json:"storyboards"`
	Microformat struct {
		PlayerMicroformatRenderer struct {
			Thumbnail struct {
				Thumbnails []struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"thumbnails"`
			} `json:"thumbnail"`
			Embed struct {
				IframeURL      string `json:"iframeUrl"`
				FlashURL       string `json:"flashUrl"`
				Width          int    `json:"width"`
				Height         int    `json:"height"`
				FlashSecureURL string `json:"flashSecureUrl"`
			} `json:"embed"`
			Title struct {
				SimpleText string `json:"simpleText"`
			} `json:"title"`
			Description struct {
				SimpleText string `json:"simpleText"`
			} `json:"description"`
			LengthSeconds      string   `json:"lengthSeconds"`
			OwnerProfileURL    string   `json:"ownerProfileUrl"`
			ExternalChannelID  string   `json:"externalChannelId"`
			AvailableCountries []string `json:"availableCountries"`
			IsUnlisted         bool     `json:"isUnlisted"`
			HasYpcMetadata     bool     `json:"hasYpcMetadata"`
			ViewCount          string   `json:"viewCount"`
			Category           string   `json:"category"`
			PublishDate        string   `json:"publishDate"`
			OwnerChannelName   string   `json:"ownerChannelName"`
			UploadDate         string   `json:"uploadDate"`
		} `json:"playerMicroformatRenderer"`
	} `json:"microformat"`
}
