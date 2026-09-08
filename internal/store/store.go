package store

import (
	"database/sql"
	"encoding/json"
	"os"

	"github.com/espcaa/hammock/internal/slack"
	"github.com/espcaa/hammock/internal/store/db"
	"github.com/zalando/go-keyring"
)

type Paths struct {
	ConfigFile string
	CacheDb    string
}

type ConfigFile struct {
	WorkspaceIDs []string `json:"workspace_ids"`
}

type Store struct {
	SlackSession *slack.SlackSession
	Paths        Paths
	db           *sql.DB
	dbq          *db.Queries
}

func New(paths Paths) *Store {
	return &Store{
		Paths: paths,
	}
}

func DefaultPaths() Paths {
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	return Paths{
		ConfigFile: configDir + "/hammock/config.json",
		CacheDb:    configDir + "/hammock/cache.db",
	}
}

func (s *Store) LoadSession() error {
	if s.SlackSession != nil {
		return nil
	}

	// to build a slack session, we need a dcookie & a workspace id arrary
	// dcookie -> keyring
	// workspace id array -> config file

	dCookie, err := keyring.Get("hammock", "dcookie")
	if err != nil {
		return nil
	}

	config, err := os.ReadFile(s.Paths.ConfigFile)
	if err != nil {
		return nil
	}

	var configData ConfigFile
	err = json.Unmarshal(config, &configData)
	if err != nil {
		return nil
	}

	s.SlackSession = &slack.SlackSession{
		DCookie:       dCookie,
		WorkspacesIds: configData.WorkspaceIDs,
	}

	return nil
}

func (s *Store) SaveSession() error {
	if s.SlackSession == nil {
		return nil
	}

	// save the dcookie to the keyring
	err := keyring.Set("hammock", "dcookie", s.SlackSession.DCookie)
	if err != nil {
		return err
	}

	// save the workspace ids to the config file
	configData := ConfigFile{
		WorkspaceIDs: s.SlackSession.WorkspacesIds,
	}
	config, err := json.Marshal(configData)
	if err != nil {
		return err
	}

	err = os.WriteFile(s.Paths.ConfigFile, config, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) ClearSession() error {
	// clear the dcookie from the keyring
	err := keyring.Delete("hammock", "dcookie")
	if err != nil {
		return err
	}

	// clear the workspace ids from the config file
	err = os.Remove(s.Paths.ConfigFile)
	if err != nil {
		return err
	}

	s.SlackSession = nil

	return nil
}
