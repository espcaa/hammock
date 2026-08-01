package slack

type Channel struct {
	ID         string `json:"id"`
	Name       string `json:"name,omitempty"`
	IsChannel  bool   `json:"is_channel"`
	IsGroup    bool   `json:"is_group"`
	IsIM       bool   `json:"is_im"`
	IsMpIM     bool   `json:"is_mpim"`
	IsPrivate  bool   `json:"is_private"`
	Created    int64  `json:"created"`
	IsArchived bool   `json:"is_archived"`
	Updated    int64  `json:"updated"`
	Creator    string `json:"creator,omitempty"`

	// userBoot fields
	IsGeneral       bool     `json:"is_general,omitempty"`
	NameNormalized  string   `json:"name_normalized,omitempty"`
	IsShared        bool     `json:"is_shared,omitempty"`
	IsFrozen        bool     `json:"is_frozen,omitempty"`
	IsOrgShared     bool     `json:"is_org_shared,omitempty"`
	IsExtShared     bool     `json:"is_ext_shared,omitempty"`
	ContextTeamID   string   `json:"context_team_id,omitempty"`
	SharedTeamIDs   []string `json:"shared_team_ids,omitempty"`
	InternalTeamIDs []string `json:"internal_team_ids,omitempty"`
	Members         []string `json:"members,omitempty"`

	// conversations.info / IM fields
	User          string   `json:"user,omitempty"`
	IsOpen        bool     `json:"is_open,omitempty"`
	LastRead      string   `json:"last_read,omitempty"`
	Priority      int      `json:"priority,omitempty"`
	UnreadCount   int      `json:"unread_count,omitempty"`
	UnreadDisplay int      `json:"unread_count_display,omitempty"`
	Latest        *Message `json:"latest,omitempty"`

	Topic struct {
		Value   string `json:"value"`
		Creator string `json:"creator"`
		LastSet int64  `json:"last_set"`
	} `json:"topic"`
	Purpose struct {
		Value   string `json:"value"`
		Creator string `json:"creator"`
		LastSet int64  `json:"last_set"`
	} `json:"purpose"`

	Properties struct {
		Tabs                []Tab `json:"tabs,omitempty"`
		PostingRestrictedTo *struct {
			Type []string `json:"type"`
			User []string `json:"user"`
		} `json:"posting_restricted_to,omitempty"`
		IsDormant bool `json:"is_dormant,omitempty"`
		Canvas    *struct {
			FileID       string `json:"file_id"`
			IsEmpty      bool   `json:"is_empty"`
			QuipThreadID string `json:"quip_thread_id"`
		} `json:"canvas,omitempty"`
	} `json:"properties"`

	PreviousNames []string `json:"previous_names,omitempty"`
}

type Tab struct {
	Type       string `json:"type"`
	Label      string `json:"label,omitempty"`
	ID         string `json:"id,omitempty"`
	IsDisabled *bool  `json:"is_disabled,omitempty"`
	Data       *struct {
		FileID           string `json:"file_id,omitempty"`
		SharedTS         string `json:"shared_ts,omitempty"`
		MuteEditUpdates  bool   `json:"mute_edit_updates,omitempty"`
		FolderBookmarkID string `json:"folder_bookmark_id,omitempty"`
	} `json:"data,omitempty"`
}
