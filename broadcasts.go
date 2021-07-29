package goinsta

import (
	"encoding/json"
	"strconv"
)

// Broadcast is live videos.
type Broadcast struct {
	insta              *Instagram
	LastLikeTs         int64
	LastCommentTs      int64
	LastCommentFetchTs int64
	LastCommentTotal   int

	ID                         int64   `json:"id"`
	MediaID                    string  `json:"media_id"`
	LivePostID                 int64   `json:"live_post_id"`
	BroadcastStatus            string  `json:"broadcast_status"`
	DashPlaybackUrl            string  `json:"dash_playback_url"`
	DashAbrPlaybackUrl         string  `json:"dash_abr_playback_url"`
	DashManifest               string  `json:"dash_manifest"`
	ExpireAt                   int64   `json:"expire_at"`
	EncodingTag                string  `json:"encoding_tag"`
	InternalOnly               bool    `json:"internal_only"`
	NumberOfQualities          int     `json:"number_of_qualities"`
	CoverFrameURL              string  `json:"cover_frame_url"`
	User                       User    `json:"broadcast_owner"`
	Cobroadcasters             []User  `json:"cobroadcasters"`
	PublishedTime              int64   `json:"published_time"`
	BroadcastMessage           string  `json:"broadcast_message"`
	OrganicTrackingToken       string  `json:"organic_tracking_token"`
	IsPlayerLiveTrace          int     `json:"is_player_live_trace_enabled"`
	IsGamingContent            int     `json:"is_gaming_content"`
	IsViewerCommentAllowed     bool    `json:"is_viewer_comment_allowed"`
	LiveCommentMentionEnabled  bool    `json:"is_live_comment_mention_enabled"`
	LiveCommmentRepliesEnabled bool    `json:"is_live_comment_replies_enabled"`
	HideFromFeedUnit           bool    `json:"hide_from_feed_unit"`
	VideoDuration              float64 `json:"video_duration"`
	Visibility                 int     `json:"visibility"`
	ResponseTs                 int64   `json:"response_timestamp"`
	Status                     string  `json:"status"`
	Dimensions                 struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"dimensions"`
	Experiments     map[string]interface{} `json:"broadcast_experiments"`
	PayViewerConfig struct {
		PayConfig struct {
			ConsumptionSheetConfig struct {
				Description               string `json:"description"`
				PrivacyDisclaimer         string `json:"privacy_disclaimer"`
				PrivacyDisclaimerLink     string `json:"privacy_disclaimer_link"`
				PrivacyDisclaimerLinkText string `json:"privacy_disclaimer_link_text"`
			} `json:"consumption_sheet_config"`
			DigitalNonConsumableProductID int64 `json:"digital_non_consumable_product_id"`
			DigitalProductID              int64 `json:"digital_product_id"`
			PayeeID                       int64 `json:"payee_id"`
			PinnedRowConfig               struct {
				ButtonTitle string `json:"button_title"`
				Description string `json:"description"`
			} `json:"pinned_row_config"`
			TierInfos struct {
				DigitalProductID int64  `json:"digital_product_id"`
				Sku              string `json:"sku"`
				SupportTier      string `json:"support_tier"`
			} `json:"tier_infos"`
		} `json:"pay_config"`
	} `json:"user_pay_viewer_config"`
}

type BroadcastComments struct {
	CommentLikesEnabled        bool      `json:"comment_likes_enabled"`
	Comments                   []Comment `json:"comments"`
	PinnedComment              Comment   `json:"pinned_comment"`
	CommentCount               int       `json:"comment_count"`
	Caption                    Caption   `json:"caption"`
	CaptionIsEdited            bool      `json:"caption_is_edited"`
	HasMoreComments            bool      `json:"has_more_comments"`
	HasMoreHeadloadComments    bool      `json:"has_more_headload_comments"`
	MediaHeaderDisplay         string    `json:"media_header_display"`
	CanViewMorePreviewComments bool      `json:"can_view_more_preview_comments"`
	LiveSecondsPerComment      int       `json:"live_seconds_per_comment"`
	IsFirstFetch               string    `json:"is_first_fetch"`
	SystemComments             []Comment `json:"system_comments"`
	CommentMuted               int       `json:"comment_muted"`
	IsViewerCommentAllowed     bool      `json:"is_viewer_comment_allowed"`
	Status                     string    `json:"status"`
}

type BroadcastLikes struct {
	Likes            int    `json:"likes"`
	BurstLikes       int    `json:"burst_likes"`
	Likers           []User `json:"likers"`
	LikeTs           int64  `json:"like_ts"`
	Status           string `json:"status"`
	PaySupporterInfo struct {
		LikeCountByTier []struct {
			BurstLikes  int         `json:"burst_likes"`
			Likers      interface{} `json:"likers"`
			Likes       int         `json:"likes"`
			SupportTier string      `json:"support_tier"`
		} `json:"like_count_by_support_tier"`
		BurskLikes int `json:"supporter_tier_burst_likes"`
		Likes      int `json:"supporter_tier_likes"`
	} `json:"user_pay_supporter_info"`
}

type BroadcastHeartbeat struct {
	ViewerCount             float64  `json:"viewer_count"`
	BroadcastStatus         string   `json:"broadcast_status"`
	CobroadcasterIds        []string `json:"cobroadcaster_ids"`
	OffsetVideoStart        float64  `json:"offset_to_video_start"`
	RequestToJoinEnabled    int      `json:"request_to_join_enabled"`
	UserPayMaxAmountReached bool     `json:"user_pay_max_amount_reached"`
	Status                  string   `json:"status"`
}

func (br *Broadcast) GetInfo() error {
	body, _, err := br.insta.sendRequest(
		&reqOptions{
			Endpoint: urlLiveComments,
			Query: map[string]string{
				"view_expired_broadcast": "false",
			},
		})
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, br)
	return err
}

// Call every 2 seconds
func (br *Broadcast) GetComments() (*BroadcastComments, error) {
	body, _, err := br.insta.sendRequest(
		&reqOptions{
			Endpoint: urlLiveComments,
			Query: map[string]string{
				"last_comment_ts":               strconv.Itoa(int(br.LastCommentTs)),
				"join_request_last_seen_ts":     "0",
				"join_request_last_fetch_ts":    "0",
				"join_request_last_total_count": "0",
			},
		})
	if err != nil {
		return nil, err
	}

	c := &BroadcastComments{}
	err = json.Unmarshal(body, c)
	if err != nil {
		return nil, err
	}

	if c.CommentCount > 0 {
		br.LastCommentTs = c.Comments[0].CreatedAt
	}
	return c, nil
}

// Call every 6 seconds
func (br *Broadcast) GetLikes() (*BroadcastLikes, error) {
	body, _, err := br.insta.sendRequest(
		&reqOptions{
			Endpoint: urlLiveComments,
			Query: map[string]string{
				"like_ts": strconv.Itoa(int(br.LastLikeTs)),
			},
		})
	if err != nil {
		return nil, err
	}

	c := &BroadcastLikes{}
	err = json.Unmarshal(body, c)
	if err != nil {
		return nil, err
	}
	br.LastLikeTs = c.LikeTs
	return c, nil
}

// Call every 3 seconds
func (br *Broadcast) GetHeartbeat() (*BroadcastHeartbeat, error) {
	body, _, err := br.insta.sendRequest(
		&reqOptions{
			Endpoint: urlLiveComments,
			IsPost:   true,
			Query: map[string]string{
				"_uuid":                 br.insta.uuid,
				"live_with_eligibility": "2", // What is this?
			},
		})
	if err != nil {
		return nil, err
	}

	c := &BroadcastHeartbeat{}
	err = json.Unmarshal(body, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (br *Broadcast) SetInsta(insta *Instagram) {
	br.insta = insta
}
