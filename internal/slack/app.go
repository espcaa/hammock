package slack

import "encoding/json"

type AppIcons struct {
	Emoji         string `json:"emoji,omitempty"`
	ImageOriginal string `json:"image_original,omitempty"`
	Image24       string `json:"image_24,omitempty"`
	Image32       string `json:"image_32,omitempty"`
	Image36       string `json:"image_36,omitempty"`
	Image48       string `json:"image_48,omitempty"`
	Image64       string `json:"image_64,omitempty"`
	Image72       string `json:"image_72,omitempty"`
	Image96       string `json:"image_96,omitempty"`
	Image128      string `json:"image_128,omitempty"`
	Image192      string `json:"image_192,omitempty"`
	Image512      string `json:"image_512,omitempty"`
	Image1024     string `json:"image_1024,omitempty"`
}

type AppCommand struct {
	Usage string `json:"usage"`
	Desc  string `json:"desc"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	App   string `json:"app"`
}

type AppScreenshot struct {
	ID            string `json:"id"`
	Image440      string `json:"image_440,omitempty"`
	Image1000     string `json:"image_1000,omitempty"`
	Image1600     string `json:"image_1600,omitempty"`
	ImageOriginal string `json:"image_original,omitempty"`
}

type AppBotUser struct {
	ID               string `json:"id"`
	Username         string `json:"username"`
	MembershipsCount int    `json:"memberships_count"`
}

type AppAuth struct {
	CreatedBy   string   `json:"created_by"`
	DateCreated string   `json:"date_created"`
	Scopes      []string `json:"scopes"`
	Username    string   `json:"username"`
	FullName    string   `json:"full_name"`
	RealName    string   `json:"real_name"`
	Icons       AppIcons `json:"icons"`
}

type BotInfo struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	AppID   string   `json:"app_id,omitempty"`
	UserID  string   `json:"user_id,omitempty"`
	Deleted bool     `json:"deleted,omitempty"`
	Icons   AppIcons `json:"icons"`
}

type AppProfile struct {
	ID                   string                     `json:"id"`
	Name                 string                     `json:"name"`
	DeveloperName        string                     `json:"developer_name"`
	Desc                 string                     `json:"desc"`
	LongDesc             string                     `json:"long_desc"`
	LongDescFormatted    string                     `json:"long_desc_formatted"`
	URL                  string                     `json:"url"`
	SupportURL           string                     `json:"support_url"`
	ConfigURL            string                     `json:"config_url"`
	AppCardColor         string                     `json:"app_card_color"`
	InstallationSummary  string                     `json:"installation_summary"`
	TeamID               string                     `json:"team_id"`
	EnterpriseID         string                     `json:"enterprise_id"`
	IsCertified          bool                       `json:"is_certified"`
	IsDirectoryPublished bool                       `json:"is_directory_published"`
	IsDistributed        bool                       `json:"is_distributed"`
	IsAIApp              bool                       `json:"is_ai_app"`
	IsAgentApp           bool                       `json:"is_agent_app"`
	IsWorkflowApp        bool                       `json:"is_workflow_app"`
	DateInstalled        int64                      `json:"date_installed"`
	Commands             map[string]AppCommand      `json:"commands"`
	Categories           map[string]json.RawMessage `json:"categories"`
	Screenshots          []AppScreenshot            `json:"screenshots"`
	Icons                AppIcons                   `json:"icons"`
	BotUser              AppBotUser                 `json:"bot_user"`
	Auth                 AppAuth                    `json:"auth"`
	SecurityCompliance   json.RawMessage            `json:"security_compliance,omitempty"`
}
