package slack

import (
	"encoding/json"
	"errors"
)

type EdgeUsersListResponse struct {
	NextMarker string        `json:"next_marker"`
	OK         bool          `json:"ok"`
	Results    []UserProfile `json:"results"`
}

type EdgeUsersListRequest struct {
	Channels     []string `json:"channels"`
	Count        int      `json:"count"`
	Filter       string   `json:"filter"`
	PresentFirst bool     `json:"present_first"`
}

func (c *Client) edgeUsersList(teamID string, requestPayload json.RawMessage) (EdgeUsersListResponse, error) {
	raw, err := c.DoEdge(teamID, "users/list", requestPayload)

	if err != nil {
		return EdgeUsersListResponse{}, err
	}

	var resp EdgeUsersListResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return EdgeUsersListResponse{}, errors.New("failed to unmarshal edge users.list response: " + err.Error() + " raw: " + string(raw))
	}
	return resp, nil
}
