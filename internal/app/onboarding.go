package app

import (
	"bytes"
	_ "embed"
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
	err := st.LoadSession()

	if err == nil && st.SlackSession != nil { // we have a  session with a dcookie, but maybe no workspaces yet
		h.loading = true
		h.loadingStatus = "Fetching tokens..."

		workspaces, err := slack.FetchTokens(st.SlackSession.DCookie)
		if err != nil {
			log.Printf("failed to fetch tokens: %v", err)
			h.loading = false
			h.errorMessage = "Failed to fetch tokens. Please try again."
			return
		}

		workspaceIds := make([]string, 0, len(workspaces))
		for id := range workspaces {
			workspaceIds = append(workspaceIds, id)
		}

		st.SlackSession.WorkspacesIds = workspaceIds
		st.SlackSession.WorkspaceSessions = workspaces

		log.Printf("Fetched workspaces: %v", workspaceIds)

		h.loadingStatus = "Saving session..."
		err = h.store.SaveSession()
		if err != nil {
			log.Printf("failed to save session: %v", err)
			h.loading = false
			h.errorMessage = "Failed to save session. Please try again."
			return
		}

		err = st.OpenCache()
		if err != nil {
			h.loading = false
			h.errorMessage = "Failed to open cache. Please try again."
			return
		}

		client := slack.NewClient(st.SlackSession)
		if client == nil {
			h.loading = false
			h.errorMessage = "Failed to create slack client. Please try again."
			return
		}

		// userboot for each workspace
		for _, workspaceId := range workspaceIds {
			h.loadingStatus = "Booting workspace: " + workspaceId

			minChannelUpdated, err := st.MinChannelUpdated(workspaceId)

			if err != nil {
				minChannelUpdated = int64(0)
				log.Printf("failed to get min channel updated for workspace %s: %v", workspaceId, err)
			} else {
				log.Printf("min channel updated for workspace %s: %d", workspaceId, minChannelUpdated)
			}

			var userbootResp *slack.UserbootResponse

			userbootResp, err = client.UserBoot(workspaceId, minChannelUpdated)
			if err != nil {
				log.Printf("failed to userboot workspace %s: %v", workspaceId, err)
				h.loading = false
				h.errorMessage = "Failed to boot workspace: " + workspaceId + ". Please try again."
				return
			}
			if userbootResp.OK {
				log.Printf("Successfully booted workspace %s", workspaceId)

				// save channels to cache
				err = st.UpsertChannels(userbootResp.Channels, workspaceId)
				if err != nil {
					log.Printf("failed to save channels for workspace %s: %v", workspaceId, err)
					h.loading = false
					h.errorMessage = "Failed to save channels for workspace: " + workspaceId + ". Please try again."
					return
				} else {
					log.Printf("Successfully saved channels for workspace %s", workspaceId)
				}
			} else {
				log.Printf("Failed to boot workspace %s, try again", workspaceId)
				h.loading = false
				h.errorMessage = "Failed to boot workspace: " + workspaceId + ". Please try again."
				return
			}
		}

		// finally navigate to the main screen
		h.mu.Lock()
		h.pendingNav = NewMainScreen(h.theme, h.router, st, client)
		h.mu.Unlock()
		h.router.Invalidate()

		return
	}
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
			// get the exec path of the current binary
			execPath, err := os.Executable()
			if err != nil {
				log.Printf("failed to get executable path: %v", err)
				return
			}

			// run the exec with the __login flag
			cmd := exec.Command(execPath, "__login")
			output, err := cmd.Output()
			if err != nil {
				log.Printf("failed to run login process: %v", err)
				a.loading = false
				return
			}

			parts := bytes.SplitN(output, []byte("|"), 2)
			if len(parts) != 2 {
				log.Printf("invalid output from login process: %s", output)
				a.errorMessage = "Invalid output from login process. Please try again."
				a.loading = false
				return
			}

			magicCode := string(parts[0])
			workspaceId := string(parts[1])
			log.Printf("Received magic code: %s", magicCode)
			log.Printf("Received workspace id: %s", workspaceId)

			a.loadingStatus = "Redeeming auth cookies..."

			dcookie, err := slack.RedeemAuthCookies(magicCode, workspaceId, nil)

			if err != nil {
				log.Printf("failed to redeem auth cookies: %v", err)
				a.loading = false
				a.errorMessage = "Failed to redeem auth cookies. Please try again."
				return
			}

			log.Printf("Received auth cookies: %v", dcookie)

			a.store.SlackSession = &slack.SlackSession{
				DCookie: dcookie,
			}

			// continue logging in

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
