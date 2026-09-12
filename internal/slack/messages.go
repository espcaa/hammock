package slack

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type MessagesResponse struct {
	Messages   []Message `json:"messages"`
	HasMore    bool      `json:"has_more"`
	NextCursor string    `json:"next_cursor"`
}

type Message struct {
	User        string          `json:"user"`
	Text        string          `json:"text"`
	Ts          string          `json:"ts"`
	Type        string          `json:"type"`
	Subtype     string          `json:"subtype,omitempty"`
	Team        string          `json:"team,omitempty"`
	ThreadTs    string          `json:"thread_ts,omitempty"`
	ReplyCount  int             `json:"reply_count,omitempty"`
	LatestReply string          `json:"latest_reply,omitempty"`
	ReplyUsers  []string        `json:"reply_users,omitempty"`
	Blocks      json.RawMessage `json:"blocks,omitempty"`
	Edited      json.RawMessage `json:"edited,omitempty"`
	BotProfile  json.RawMessage `json:"bot_profile,omitempty"`
	Files       []File          `json:"files,omitempty"`
	Raw         json.RawMessage `json:"-"`
	BotID       string          `json:"bot_id,omitempty"`
	AppID       string          `json:"app_id,omitempty"`
	Username    string          `json:"username,omitempty"`
	Icons       *AppIcons       `json:"icons,omitempty"`
}

type File struct {
	Id                 string      `json:"id"`
	Created            int64       `json:"created"`
	Timestamp          json.Number `json:"timestamp"`
	MimeType           string      `json:"mimetype"`
	Filetype           string      `json:"filetype"`
	PrettyType         string      `json:"pretty_type"`
	User               string      `json:"user"`
	UserTeam           string      `json:"user_team"`
	Editable           bool        `json:"editable"`
	Size               int64       `json:"size"`
	Mode               string      `json:"mode"`
	IsExternal         bool        `json:"is_external"`
	ExternalType       string      `json:"external_type"`
	IsPublic           bool        `json:"is_public"`
	PublicURLShared    bool        `json:"public_url_shared"`
	DisplayAsBot       bool        `json:"display_as_bot"`
	Username           string      `json:"username,omitempty"`
	Name               string      `json:"name,omitempty"`
	Title              string      `json:"title,omitempty"`
	UrlPrivate         string      `json:"url_private,omitempty"`
	UrlPrivateDownload string      `json:"url_private_download,omitempty"`
	MediaDisplayType   string      `json:"media_display_type,omitempty"`

	// very smol
	Thumb64  string `json:"thumb_64,omitempty"`
	Thumb80  string `json:"thumb_80,omitempty"`
	Thumb160 string `json:"thumb_160,omitempty"`

	Thumb360  string `json:"thumb_360,omitempty"`
	Thumb360W int    `json:"thumb_360_w,omitempty"`
	Thumb360H int    `json:"thumb_360_h,omitempty"`
	Thumb480  string `json:"thumb_480,omitempty"`
	Thumb480W int    `json:"thumb_480_w,omitempty"`
	Thumb480H int    `json:"thumb_480_h,omitempty"`

	Thumb720   string `json:"thumb_720,omitempty"`
	Thumb720W  int    `json:"thumb_720_w,omitempty"`
	Thumb720H  int    `json:"thumb_720_h,omitempty"`
	Thumb800   string `json:"thumb_800,omitempty"`
	Thumb800W  int    `json:"thumb_800_w,omitempty"`
	Thumb800H  int    `json:"thumb_800_h,omitempty"`
	Thumb960   string `json:"thumb_960,omitempty"`
	Thumb960W  int    `json:"thumb_960_w,omitempty"`
	Thumb960H  int    `json:"thumb_960_h,omitempty"`
	Thumb1024  string `json:"thumb_1024,omitempty"`
	Thumb1024W int    `json:"thumb_1024_w,omitempty"`
	Thumb1024H int    `json:"thumb_1024_h,omitempty"`

	ThumbTiny  string `json:"thumb_tiny,omitempty"`
	ThumbPdf   string `json:"thumb_pdf,omitempty"`
	ThumbVideo string `json:"thumb_video,omitempty"`
	ThumbGif   string `json:"thumb_360_gif,omitempty"`

	OriginalW int `json:"original_w,omitempty"`
	OriginalH int `json:"original_h,omitempty"`

	Permalink       string `json:"permalink,omitempty"`
	PermalinkPublic string `json:"permalink_public,omitempty"`
	HasRichPreview  bool   `json:"has_rich_preview,omitempty"`
	IsStarred       bool   `json:"is_starred,omitempty"`
	FileAccess      string `json:"file_access,omitempty"`
}

type ConversationHistoryResponse struct {
	OK               bool      `json:"ok"`
	Messages         []Message `json:"messages"`
	ResponseMetadata struct {
		NextCursor string `json:"next_cursor"`
	} `json:"response_metadata"`
}

func (c *Client) GetConversationHistory(teamId, channelID string, latest string, limit int) ([]Message, error) {
	params := url.Values{}
	params.Set("channel", channelID)
	if latest != "" {
		params.Set("latest", latest)
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	} else {
		params.Set("limit", "28")
	}
	params.Set("no_user_profile", "true")

	raw, err := c.DoWithQuery(teamId, "conversations.history", params, nil)
	if err != nil {
		return nil, err
	}

	var resp ConversationHistoryResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	return resp.Messages, nil
}
