package store

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/espcaa/hammock/internal/slack"
	"github.com/espcaa/hammock/internal/store/db"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaSQL string

func (s *Store) OpenCache() error {
	os.MkdirAll(filepath.Dir(s.Paths.CacheDb), 0o755)
	database, err := sql.Open("sqlite", s.Paths.CacheDb+"?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
	if err != nil {
		return err
	}
	if _, err := database.Exec(schemaSQL); err != nil {
		return err
	}
	s.db, s.dbq = database, db.New(database)
	return nil
}

func (s *Store) MinChannelUpdated(teamId string) (int64, error) {
	ctx := context.Background()

	max, err := s.dbq.MinChannelUpdated(ctx, teamId)
	if err != nil {
		return 0, err
	}
	return max, nil
}

func (s *Store) ListChannels(teamId, channelType string) ([]db.Channel, error) {
	ctx := context.Background()

	channels, err := s.dbq.ListChannels(ctx, db.ListChannelsParams{
		TeamID: teamId,
		Type:   channelType,
	})
	if err != nil {
		return nil, err
	}
	return channels, nil
}

func (s *Store) UpsertChannels(channels []slack.Channel, teamId string) error {
	ctx := context.Background()

	for _, channel := range channels {

		jsonMemberData, err := json.Marshal(channel.Members)
		if err != nil {
			return err
		}

		jsonTopicData, err := json.Marshal(channel.Topic)
		if err != nil {
			return err
		}

		channelType := "channel-public"
		if channel.IsPrivate {
			channelType = "channel-private"
		} else if channel.IsIM {
			channelType = "im"
		}

		err = s.dbq.UpsertChannel(ctx, db.UpsertChannelParams{
			TeamID:  teamId,
			ID:      channel.ID,
			Name:    channel.Name,
			Type:    channelType,
			Unread:  int64(channel.UnreadCount),
			Updated: channel.Updated,
			Members: jsonMemberData,
			Topic:   jsonTopicData,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
