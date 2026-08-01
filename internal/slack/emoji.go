package slack

type Emoji struct {
	Name    string `json:"name"`
	Url     string `json:"value"`
	Unicode string `json:"unicode,omitempty"`
	Updated int64  `json:"updated"`
}
