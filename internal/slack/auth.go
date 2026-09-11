package slack

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

type slackAuthResponse struct {
	OK           bool                        `json:"ok"`
	TokenResults map[string]slackTokenResult `json:"token_results"`
}

type SlackSession struct {
	DCookie           string   `json:"d_cookie"`
	WorkspacesIds     []string `json:"workspaces_ids"`
	WorkspaceSessions map[string]WorkspaceSession
}

type WorkspaceSession struct {
	Token        string `json:"token"`
	UserID       string `json:"user_id"`
	TeamName     string `json:"team_name"`
	TeamURL      string `json:"team_url"`
	TeamIcon     string `json:"team_icon"`
	EnterpriseID string `json:"enterprise_id,omitempty"`
	Domain       string `json:"domain,omitempty"`
}

type slackTokenResult struct {
	OK   bool   `json:"ok"`
	User string `json:"user"`
	Team struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"team"`
}

type authConfigTeam struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	URL                string `json:"url"`
	Token              string `json:"token"`
	UserID             string `json:"user_id"`
	EnterpriseID       string `json:"enterprise_id"`
	EnterpriseAPIToken string `json:"enterprise_api_token"`
	Icon               struct {
		Image68 string `json:"image_68"`
	} `json:"icon"`
}

type authConfigResponse struct {
	Teams map[string]authConfigTeam `json:"teams"`
}

func RedeemAuthCookies(magicToken string, workspaceId string, cookies []Cookie) (string, error) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	endpoint := fmt.Sprintf("https://app.slack.com/api/auth.loginMagicBulk?magic_tokens=z-app-%s-%s&ssb=1", workspaceId, magicToken)
	targetURL, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	var httpCookies []*http.Cookie
	for _, c := range cookies {
		httpCookies = append(httpCookies, &http.Cookie{
			Name:  c.Name,
			Value: c.Value,
		})
	}

	client.Jar.SetCookies(targetURL, httpCookies)

	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "d" {
			return cookie.Value, nil
		}
	}
	return "", errors.New("no d cookie found")
}

func FetchTokens(dCookie string) (map[string]WorkspaceSession, error) {
	// make a request to slack & get the tokens from the auth config json embedded in the html
	endpoint := "https://app.slack.com/auth?app=client&return_to=%2Fclient&teams=&iframe=1"
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Cookie", "d="+dCookie)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	prefix := "JSON.stringify("
	suffix := ");"
	start := strings.Index(bodyStr, prefix)
	if start == -1 {
		return nil, errors.New("could not find config JSON in auth response")
	}
	start += len(prefix)
	end := strings.Index(bodyStr[start:], suffix)
	if end == -1 {
		return nil, errors.New("could not find end of config json")
	}
	configJSON := bodyStr[start : start+end]

	var config authConfigResponse
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("failed to parse auth config: %w", err)
	}

	workspaces := make(map[string]WorkspaceSession)

	// Collect enterprise teams first for URL/token lookup
	enterpriseTeams := make(map[string]authConfigTeam)
	for id, team := range config.Teams {
		if strings.HasPrefix(id, "E") {
			enterpriseTeams[id] = team
		}
	}

	for id, team := range config.Teams {
		if strings.HasPrefix(id, "E") {
			continue
		}

		token := team.Token
		teamURL := team.URL

		// Use enterprise URL/token if this workspace belongs to an enterprise
		if team.EnterpriseID != "" {
			if ent, ok := enterpriseTeams[team.EnterpriseID]; ok {
				teamURL = ent.URL
				if ent.Token != "" {
					token = ent.Token
				}
			}
		}
		if team.EnterpriseAPIToken != "" {
			token = team.EnterpriseAPIToken
		}

		workspaces[id] = WorkspaceSession{
			Token:        token,
			UserID:       team.UserID,
			TeamName:     team.Name,
			TeamURL:      teamURL,
			TeamIcon:     team.Icon.Image68,
			EnterpriseID: team.EnterpriseID,
		}
	}

	if len(workspaces) == 0 {
		return nil, errors.New("no workspaces found in auth config")
	}

	return workspaces, nil
}

func GenerateSSBParams() string {
	instanceID := uuid.New().String()

	h := sha256.New()
	h.Write([]byte("hammock-" + instanceID))
	vid := fmt.Sprintf("%x", h.Sum(nil))[:32]

	version := "4.48.102"

	return fmt.Sprintf("ssb_vid=%s&ssb_instance_id=%s&v=%s&from_desktop_app=1", vid, instanceID, version)
}
