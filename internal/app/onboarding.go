package app

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"os"
	"os/exec"
	"sync"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/espcaa/hammock/internal/slack"
	"github.com/espcaa/hammock/internal/store"
	"github.com/espcaa/hammock/internal/ui"
)

//go:embed assets/onboarding.jpg
var onboarding []byte

type OnboardingScreen struct {
	theme         *ui.Theme
	router        *Router
	store         *store.Store
	loginButton   ui.Button
	onboardImg    ui.Image
	loading       bool
	loadingStatus string
	errorMessage  string
	pendingNav    Screen
	mu            sync.Mutex
}

func FinishLogin(h *OnboardingScreen, st *store.Store) {
	if err := st.LoadSession(); err != nil {
		fail(h, "Failed to load session", err)
		return
	}
	if st.SlackSession == nil {
		return // no saved session yet; wait for the user to log in
	}

	h.loading = true
	h.loadingStatus = "Fetching tokens..."

	workspaces, err := slack.FetchTokens(st.SlackSession.DCookie)
	if err != nil {
		fail(h, "Failed to fetch tokens", err)
		return
	}

	st.SlackSession.WorkspacesIds = mapKeys(workspaces)
	st.SlackSession.WorkspaceSessions = workspaces
	log.Printf("Fetched workspaces: %v", st.SlackSession.WorkspacesIds)

	h.loadingStatus = "Saving session..."
	if err := st.SaveSession(); err != nil {
		fail(h, "Failed to save session", err)
		return
	}

	if err := st.OpenCache(); err != nil {
		fail(h, "Failed to open cache", err)
		return
	}

	client := slack.NewClient(st.SlackSession)
	if client == nil {
		fail(h, "Failed to create slack client", nil)
		return
	}

	for _, workspaceID := range st.SlackSession.WorkspacesIds {
		h.loadingStatus = "Booting workspace: " + workspaceID
		if err := bootWorkspace(client, st, workspaceID); err != nil {
			fail(h, fmt.Sprintf("Failed to %s", err), nil)
			return
		}
		log.Printf("Booted workspace %s", workspaceID)
	}

	h.mu.Lock()
	h.pendingNav = NewMainScreen(h.theme, h.router, st, client)
	h.mu.Unlock()
	h.router.Invalidate()
}

// bootWorkspace fetches the client.init profile, channel list and DMs for a
// single workspace, then caches them.
func bootWorkspace(client *slack.Client, st *store.Store, workspaceID string) error {
	min := int64(0)
	if v, err := st.MinChannelUpdated(workspaceID); err != nil {
		log.Printf("min channel updated for %s unavailable: %v", workspaceID, err)
	} else {
		min = v
	}

	initResp, err := client.ClientInit(workspaceID, min)
	if err != nil {
		return fmt.Errorf("init client for workspace %s: %w", workspaceID, err)
	}

	selfRaw, err := json.Marshal(initResp.Self)
	if err != nil {
		return fmt.Errorf("encode self for workspace %s: %w", workspaceID, err)
	}
	if err := st.SaveSelf(workspaceID, selfRaw); err != nil {
		return fmt.Errorf("save self for workspace %s: %w", workspaceID, err)
	}
	mergeWorkspace(st, workspaceID, initResp)

	channelResp, err := client.GetChannels(workspaceID, min)
	if err != nil {
		return fmt.Errorf("boot workspace %s: %w", workspaceID, err)
	}
	if !channelResp.OK {
		return fmt.Errorf("slack refused to boot workspace %s", workspaceID)
	}

	if err := st.UpsertChannels(channelResp.Channels.Channels, workspaceID); err != nil {
		return fmt.Errorf("save channels for workspace %s: %w", workspaceID, err)
	}
	if err := st.UpsertIms(channelResp.Channels.Ims, workspaceID); err != nil {
		return fmt.Errorf("save DMs for workspace %s: %w", workspaceID, err)
	}
	return nil
}

// mergeWorkspace copies the richer workspace metadata from client.init into
// the in-memory session.
func mergeWorkspace(st *store.Store, workspaceID string, initResp *slack.ClientInitResponse) {
	ws, ok := st.SlackSession.WorkspaceSessions[workspaceID]
	if !ok {
		return
	}
	for _, w := range initResp.Workspaces {
		if w.ID == workspaceID {
			ws.TeamName, ws.Domain = w.Name, w.Domain
			ws.TeamIcon = w.Icon.Image132
			break
		}
	}
	st.SlackSession.WorkspaceSessions[workspaceID] = ws
}

func mapKeys(m map[string]slack.WorkspaceSession) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// fail stops the loading spinner, shows a user-facing message and logs the
// underlying error once.
func fail(h *OnboardingScreen, msg string, err error) {
	if err != nil {
		log.Printf("%s: %v", msg, err)
	} else {
		log.Printf("%s", msg)
	}
	h.loading = false
	h.errorMessage = msg + " Please try again."
}

func NewOnboardingScreen(th *ui.Theme, r *Router, st *store.Store) *OnboardingScreen {
	img, _, err := image.Decode(bytes.NewReader(onboarding))
	if err != nil {
		log.Fatalf("decode onboarding.jpg: %v", err)
	}

	oi := th.Image(img, 300, 0)
	oi.Fill = true

	h := &OnboardingScreen{
		theme:      th,
		router:     r,
		store:      st,
		onboardImg: oi,
	}

	// one-time startup auth check

	go func() {
		FinishLogin(h, st)
	}()

	return h
}

func (a *OnboardingScreen) Layout(gtx layout.Context) layout.Dimensions {
	// check & apply pending navigation requests
	a.mu.Lock()
	next := a.pendingNav
	a.pendingNav = nil
	a.mu.Unlock()
	if next != nil {
		a.router.Push(gtx, next)
	}

	// handle button click for login

	if a.loginButton.Click.Clicked(gtx) {
		a.loading = true
		a.loadingStatus = "Login into slack on the webview!"

		// launch goroutine to run the login process
		go func() {
			execPath, err := os.Executable()
			if err != nil {
				fail(a, "Failed to locate app binary", err)
				return
			}

			cmd := exec.Command(execPath, "__login")
			output, err := cmd.Output()
			if err != nil {
				fail(a, "Login process failed", err)
				return
			}

			parts := bytes.SplitN(output, []byte("|"), 2)
			if len(parts) != 2 {
				fail(a, "Invalid output from login process", nil)
				return
			}

			magicCode := string(parts[0])
			workspaceId := string(parts[1])

			a.loadingStatus = "Redeeming auth cookies..."

			dcookie, err := slack.RedeemAuthCookies(magicCode, workspaceId, nil)
			if err != nil {
				fail(a, "Failed to redeem auth cookies", err)
				return
			}

			a.store.SlackSession = &slack.SlackSession{
				DCookie: dcookie,
			}

			FinishLogin(a, a.store)
		}()
	}

	// draw the actual screen

	paint.Fill(gtx.Ops, a.theme.Bg)

	return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.onboardImg.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// loading screen
			if a.loading {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return a.theme.Spinner(unit.Dp(64)).Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return a.theme.Label(a.loadingStatus, unit.Sp(15), font.Thin, false).Layout(gtx)
						}),
					)
				})
			} else {
				// normal screen
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return a.theme.Label("Welcome to Hammock!", unit.Sp(30), font.Bold, false).Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
						layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							btn := a.theme.Button(&a.loginButton, "Login")
							btn.TextSize = unit.Sp(20)
							btn.TextWeight = font.Medium
							return btn.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if a.errorMessage != "" {
								text := a.theme.Label(a.errorMessage, unit.Sp(15), font.Thin, false)
								text.Color = a.theme.Danger
								return text.Layout(gtx)
							}
							return layout.Dimensions{}
						}),
					)
				})
			}
		}),
	)
}

func (h *OnboardingScreen) WindowOptions() []app.Option {
	return []app.Option{
		app.Size(unit.Dp(300), unit.Dp(400)),
		app.MaxSize(unit.Dp(300), unit.Dp(400)),
		app.MinSize(unit.Dp(300), unit.Dp(400)),
	}
}
