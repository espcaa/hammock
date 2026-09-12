package slack

import (
	"encoding/json"
)

type UserProfile struct {
	ID             string `json:"id"`
	Color          string `json:"color"`
	IsBot          bool   `json:"is_bot"`
	Timezone       string `json:"tz"`
	TimezoneLabel  string `json:"tz_label"`
	TimezoneOffset int64  `json:"tz_offset"`
	Name           string `json:"name"`
	RealName       string `json:"real_name"`
	Profile        struct {
		DisplayName string `json:"display_name"`
		RealName    string `json:"real_name"`
		AvatarHash  string `json:"avatar_hash"`
		Fields      map[string]struct {
			Value string `json:"value"`
			Alt   string `json:"alt"`
		} `json:"fields"`
		Title       string `json:"title"`
		Phone       string `json:"phone"`
		StatusText  string `json:"status_text"`
		StatusEmoji string `json:"status_emoji"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
	} `json:"profile"`
}

func (c *Client) GetUsersFromChannel(teamID string, channelID string, limit int) ([]UserProfile, error) {
	request := EdgeUsersListRequest{
		Channels:     []string{channelID},
		Count:        limit,
		Filter:       "everyone AND NOT bots AND NOT apps",
		PresentFirst: true,
	}

	requestPayload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := c.edgeUsersList(teamID, requestPayload)
	if err != nil {
		return nil, err
	}

	return resp.Results, nil
}
