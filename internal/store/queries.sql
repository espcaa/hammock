-- name: UpsertChannel :exec
INSERT INTO channels (team_id, id, name, type, unread, mentions, updated, members, topic)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT (team_id, id) DO UPDATE SET
  name = excluded.name, type = excluded.type,
  unread = excluded.unread, updated = excluded.updated,
  members = excluded.members, topic = excluded.topic;

-- name: ListChannels :many
SELECT * FROM channels
WHERE team_id = ? AND type = ? ORDER BY name;

-- name: UpsertIm :exec
INSERT INTO ims (team_id, id, user, unreads, updated)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT (team_id, id) DO UPDATE SET
  user = excluded.user, unreads = excluded.unreads, updated = excluded.updated;

-- name: ListIms :many
SELECT * FROM ims
WHERE team_id = ? ORDER BY user;

-- name: GetMinChannelUpdated :one
SELECT CAST(COALESCE(MAX(updated), 0) AS INTEGER) FROM channels WHERE team_id = ?;
