-- +migrate Up
CREATE TABLE IF NOT EXISTS hatena_bookmark (
  "id" SERIAL NOT NULL PRIMARY KEY,
  "example_id" SERIAL NOT NULL,
  "title" TEXT NOT NULL,
  "screenshot" TEXT NOT NULL,
  "entry_url" TEXT NOT NULL,
  "count" INT NOT NULL,
  "url" TEXT NOT NULL,
  "eid" TEXT NOT NULL,
  CONSTRAINT hatena_bookmark_example_id_fkey FOREIGN KEY ("example_id") REFERENCES example("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS "example_id_idx_hatena_bookmark" ON hatena_bookmark ("example_id");
CREATE UNIQUE INDEX IF NOT EXISTS "url_idx_hatena_bookmark" ON hatena_bookmark ("url");

CREATE TABLE IF NOT EXISTS bookmark (
  "hatena_bookmark_id" SERIAL NOT NULL,
  "user" TEXT NOT NULL,
  "comment" TEXT NOT NULL,
  "timestamp" timestamp NOT NULL,
  CONSTRAINT bookmark_hatena_bookmark_id_fkey FOREIGN KEY ("hatena_bookmark_id") REFERENCES hatena_bookmark("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS "hatena_bookmark_id_user_idx_bookmark" ON bookmark ("hatena_bookmark_id", "user");

-- +migrate Down
DROP INDEX "hatena_bookmark_id_user_idx_bookmark";
DROP INDEX "example_id_idx_hatena_bookmark";
DROP INDEX "url_idx_hatena_bookmark";

DROP TABLE bookmark;
DROP TABLE hatena_bookmark;
