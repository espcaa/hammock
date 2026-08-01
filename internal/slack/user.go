package slack

type UserProfile struct {
	ID             string `json:"id"`
	Color          string `json:"color"`
	IsBot          bool   `json:"is_bot"`
	Timezone       string `json:"tz"`
	TimezoneLabel  string `json:"tz_label"`
	TimezoneOffset int64  `json:"tz_offset"`
	Profile        struct {
		DisplayName string `json:"display_name"`
		RealName    string `json:"real_name"`
		AvatarHash  string `json:"avatar_hash"`
		Title       string `json:"title"`
		Phone       string `json:"phone"`
		StatusText  string `json:"status_text"`
		StatusEmoji string `json:"status_emoji"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
	} `json:"profile"`
}
