package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/espcaa/hammock/internal/slack"
	"github.com/espcaa/hammock/internal/store/db"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaSQL string

const schemaHashKey = "schema_hash"

// messageCacheLimit caps how many cached messages per channel ListMessages
// returns (newest first). Safer to over-fetch here than to silently truncate
// history.
const messageCacheLimit = 1000

func (s *Store) OpenCache() error {
	os.MkdirAll(filepath.Dir(s.Paths.CacheDb), 0o755)
	database, err := sql.Open("sqlite", s.Paths.CacheDb+"?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
	if err != nil {
		return err
	}
	if err := s.applySchema(database); err != nil {
		database.Close()
		return err
	}
	s.db, s.dbq = database, db.New(database)
	return nil
}

func (s *Store) applySchema(d *sql.DB) error {
	if _, err := d.Exec("CREATE TABLE IF NOT EXISTS meta (key TEXT PRIMARY KEY, value TEXT NOT NULL)"); err != nil {
		return err
	}

	var stored string
	err := d.QueryRow("SELECT value FROM meta WHERE key = ?", schemaHashKey).Scan(&stored)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == nil && stored == schemaHash() {
		return nil
	}

	if err := dropAllTables(d); err != nil {
		return err
	}
	if _, err := d.Exec(schemaSQL); err != nil {
		return err
	}
	_, err = d.Exec("INSERT OR REPLACE INTO meta (key, value) VALUES (?, ?)", schemaHashKey, schemaHash())
	return err
}

func schemaHash() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(schemaSQL)))
}

func dropAllTables(d *sql.DB) error {
	rows, err := d.Query("SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%'")
	if err != nil {
		return err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}
		tables = append(tables, name)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for _, name := range tables {
		if _, err := d.Exec("DROP TABLE " + strconv.Quote(name)); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) MinChannelUpdated(teamId string) (int64, error) {
	ctx := context.Background()

	max, err := s.dbq.GetMinChannelUpdated(ctx, teamId)
	if err != nil {
		return 0, err
	}
	return max, nil
}

func (s *Store) ListChannels(teamId string, channelTypes []string) ([]db.Channel, error) {
	ctx := context.Background()

	var channels = []db.Channel{}
	for _, channelType := range channelTypes {
		channelsOfType, err := s.dbq.ListChannels(ctx, db.ListChannelsParams{
			TeamID: teamId,
			Type:   channelType,
		})
		if err != nil {
			return nil, err
		}
		channels = append(channels, channelsOfType...)
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
			TeamID:   teamId,
			ID:       channel.ID,
			Name:     channel.Name,
			Type:     channelType,
			Unread:   int64(0),
			Mentions: int64(0),
			Updated:  channel.Updated,
			Members:  jsonMemberData,
			Topic:    jsonTopicData,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) ListIms(teamId string) ([]db.Im, error) {
	ctx := context.Background()

	ims, err := s.dbq.ListIms(ctx, teamId)
	if err != nil {
		return nil, err
	}
	return ims, nil
}

func (s *Store) UpsertIms(ims []slack.Im, teamId string) error {
	ctx := context.Background()

	for _, im := range ims {
		err := s.dbq.UpsertIm(ctx, db.UpsertImParams{
			TeamID:  teamId,
			ID:      im.ID,
			User:    im.User,
			Unreads: int64(0), // TODO: implement unreads
			Updated: im.Updated,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) SaveSelf(teamID string, self json.RawMessage) error {
	ctx := context.Background()
	return s.dbq.UpsertMeta(ctx, db.UpsertMetaParams{
		Key:   fmt.Sprintf("self:%s", teamID),
		Value: string(self),
	})
}

func (s *Store) LoadSelf(teamID string) (json.RawMessage, error) {
	ctx := context.Background()
	v, err := s.dbq.GetMeta(ctx, "self:"+teamID)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(v), nil
}

func (s *Store) UpsertMessages(messages []slack.Message, teamId string, channelId string) error {
	ctx := context.Background()

	for _, message := range messages {
		jsonBlocks, err := json.Marshal(message.Blocks)
		if err != nil {
			return err
		}

		jsonRaw, err := json.Marshal(message.Raw)
		if err != nil {
			return err
		}

		err = s.dbq.UpsertMessage(ctx, db.UpsertMessageParams{
			TeamID:    teamId,
			ChannelID: channelId,
			Ts:        message.Ts,
			User:      sql.NullString{String: message.User, Valid: message.User != ""},
			Text:      sql.NullString{String: message.Text, Valid: message.Text != ""},
			ThreadTs:  sql.NullString{String: message.ThreadTs, Valid: message.ThreadTs != ""},
			Blocks:    jsonBlocks,
			Raw:       jsonRaw,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) ListMessages(teamId string, channelId string) ([]db.Message, error) {
	ctx := context.Background()

	messages, err := s.dbq.ListMessages(ctx, db.ListMessagesParams{
		TeamID:    teamId,
		ChannelID: channelId,
		Limit:     messageCacheLimit,
	})
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *Store) UpsertUsers(users []slack.UserProfile, teamId string) error {
	ctx := context.Background()

	for _, user := range users {
		jsonProfile, err := json.Marshal(user.Profile)
		if err != nil {
			return err
		}

		err = s.dbq.UpsertUsers(ctx, db.UpsertUsersParams{
			TeamID:   teamId,
			ID:       user.ID,
			Color:    sql.NullString{String: user.Color, Valid: user.Color != ""},
			Name:     user.Name,
			RealName: sql.NullString{String: user.RealName, Valid: user.RealName != ""},
			IsBot: int64(func() int {
				if user.IsBot {
					return 1
				}
				return 0
			}()),
			Timezone:       sql.NullString{String: user.Timezone, Valid: user.Timezone != ""},
			TimezoneOffset: sql.NullInt64{Int64: int64(user.TimezoneOffset), Valid: true},
			Profile:        jsonProfile,
			Updated:        time.Now().Unix(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) ListUsers(teamId string) ([]db.User, error) {
	ctx := context.Background()

	users, err := s.dbq.ListUsers(ctx, teamId)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Store) GetUser(teamId string, userId string) (db.User, error) {
	ctx := context.Background()

	user, err := s.dbq.GetUser(ctx, db.GetUserParams{
		TeamID: teamId,
		ID:     userId,
	})
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}
