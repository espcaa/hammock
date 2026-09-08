CREATE TABLE IF NOT EXISTS meta (
  key   TEXT PRIMARY KEY,
  value TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS channels (
  team_id    TEXT NOT NULL,
  id         TEXT NOT NULL,
  name       TEXT NOT NULL,
  type       TEXT NOT NULL,
  unread     INTEGER NOT NULL DEFAULT 0, -- 0 or 1 depending on if there are unread messages
  mentions   INTEGER NOT NULL DEFAULT 0, -- number of mentions in the channel
  updated    INTEGER NOT NULL,
  members    JSON,
  topic      JSON,
  PRIMARY KEY (team_id, id)
);

CREATE TABLE IF NOT EXISTS ims (
    team_id    TEXT NOT NULL,
    id         TEXT NOT NULL,
    user       TEXT NOT NULL,
    unreads     INTEGER NOT NULL DEFAULT 0, -- number of unread messages
    updated    INTEGER NOT NULL,
    PRIMARY KEY (team_id, id)
);

CREATE TABLE IF NOT EXISTS messages (
  team_id     TEXT NOT NULL,
  channel_id  TEXT NOT NULL,
  ts          TEXT NOT NULL,
  user        TEXT,
  text        TEXT,
  thread_ts   TEXT,
  blocks      JSON,
  raw         JSON NOT NULL,
  PRIMARY KEY (team_id, channel_id, ts)
);


CREATE INDEX IF NOT EXISTS idx_messages_thread ON messages (team_id, channel_id, thread_ts);
