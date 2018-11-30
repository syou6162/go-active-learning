-- +migrate Up
CREATE TABLE IF NOT EXISTS tweet (
  "id" SERIAL NOT NULL PRIMARY KEY,
  "example_id" SERIAL NOT NULL,

  "created_at" timestamp NOT NULL,
  "id_str" TEXT NOT NULL,
  "full_text" TEXT NOT NULL,
  "favorite_count" INT NOT NULL,
  "retweet_count" INT NOT NULL,
  "lang" TEXT NOT NULL,

  "screen_name" TEXT NOT NULL,
  "name" TEXT NOT NULL,
  "profile_image_url" TEXT NOT NULL,

  CONSTRAINT tweet_example_id_fkey FOREIGN KEY ("example_id") REFERENCES example("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS "example_id_idx_tweet" ON tweet ("example_id");
CREATE UNIQUE INDEX IF NOT EXISTS "example_id_id_str_idx_tweet" ON tweet ("example_id", "id_str");

-- +migrate Down
DROP INDEX "example_id_id_str_idx_tweet";
DROP INDEX "example_id_idx_tweet";
DROP TABLE tweet;
