package slack

type Cookie struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Domain string `json:"domain"`
	Path   string `json:"path"`
}

type WorkspaceSession struct {
	Token        string `json:"token"`
	UserID       string `json:"user_id"`
	TeamName     string `json:"team_name"`
	TeamURL      string `json:"team_url"`
	TeamIcon     string `json:"team_icon"`
	EnterpriseID string `json:"enterprise_id,omitempty"`
}

type SlackSession struct {
	DCookie    string                      `json:"d_cookie"`
	Workspaces map[string]WorkspaceSession `json:"workspaces"`
}

type UserbootResponse struct {
	OK                 bool   `json:"ok"`
	AppCommandsCacheTs string `json:"app_commands_cache_ts"`
	AccountType        struct {
		IsAdmin        bool `json:"is_admin"`
		IsOwner        bool `json:"is_owner"`
		IsPrimaryOwner bool `json:"is_primary_owner"`
	} `json:"account_type"`
	Channels []Channel `json:"channels"`
	Ims      []Channel `json:"ims"`
	Self     struct {
		ID                string `json:"id"`
		Name              string `json:"name"`
		IsBot             bool   `json:"is_bot"`
		Updated           int64  `json:"updated"`
		IsAppUser         bool   `json:"is_app_user"`
		Deleted           bool   `json:"deleted"`
		CompactColor      string `json:"color"` // used to display username in compact mode
		RealName          string `json:"real_name"`
		Timezone          string `json:"tz"`
		TimezoneLabel     string `json:"tz_label"`
		TimezoneOffset    int64  `json:"tz_offset"`
		IsAdmin           bool   `json:"is_admin"`
		IsOwner           bool   `json:"is_owner"`
		IsPrimaryOwner    bool   `json:"is_primary_owner"`
		IsRestricted      bool   `json:"is_restricted"`
		IsUltraRestricted bool   `json:"is_ultra_restricted"`
		FirstLogin        int64  `json:"first_login"`
		Profile           struct {
			RealName               string `json:"real_name"`
			DisplayName            string `json:"display_name"`
			AvatarHash             string `json:"avatar_hash"`
			RealNameNormalized     string `json:"real_name_normalized"`
			DisplayNameNormalized  string `json:"display_name_normalized"`
			ImageOriginal          string `json:"image_original"`
			IsCustomImage          bool   `json:"is_custom_image"`
			FirstName              string `json:"first_name"`
			LastName               string `json:"last_name"`
			Team                   string `json:"team"`
			Title                  string `json:"title"`
			Pronouns               string `json:"pronouns"`
			Phone                  string `json:"phone"`
			Skype                  string `json:"skype"`
			StatusText             string `json:"status_text"`
			StatusEmoji            string `json:"status_emoji"`
			StatusEmojiDisplayInfo []struct {
				DisplayURL string  `json:"display_url"`
				Unicode    *string `json:"unicode"`
			} `json:"status_emoji_display_info"`
			StatusExpiration   int64  `json:"status_expiration"`
			StartDate          string `json:"start_date"`
			OutOfOfficeMessage string `json:"ooo_message"`
		} `json:"profile"`
	} `json:"self"`
	Workspaces []struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Url    string `json:"url"`
		Domain string `json:"domain"`
		Icon   struct {
			ImageDefault bool   `json:"image_default"`
			Image68      string `json:"image_68"`
			Image132     string `json:"image_132"`
		} `json:"icon"`
	} `json:"workspaces"`
}
