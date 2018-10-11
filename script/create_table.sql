CREATE TABLE IF NOT EXISTS example (
  "id" SERIAL NOT NULL PRIMARY KEY,
  "url" TEXT NOT NULL,
  "label" INT NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS "url_idx_example" ON example ("url");
CREATE INDEX IF NOT EXISTS "label_updated_at_idx_example" ON example ("label", "updated_at" DESC);

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO nobody;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO nobody;
