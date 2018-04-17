CREATE TABLE IF NOT EXISTS example (
  "id" SERIAL,
  "url" TEXT NOT NULL,
  "label" INT NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS "url_idx_example" ON example ("url");
