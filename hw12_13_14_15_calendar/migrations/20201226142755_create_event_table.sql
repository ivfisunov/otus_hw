-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS events (
  id integer NOT NULL PRIMARY KEY,
  title text NOT NULL,
  date timestamp without time zone NOT NULL,
  duration bigint NOT NULL,
  desctiption text,
  user_id integer,
  notify_before bigint
  );

CREATE UNIQUE INDEX events_date_idx ON events(date);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS events;
DROP INDEX IF EXISTS events_date_idx;
