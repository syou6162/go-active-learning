CREATE TABLE IF NOT EXISTS example (
  "id" SERIAL NOT NULL PRIMARY KEY,
  "url" TEXT NOT NULL,
  "label" INT NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS "url_idx_example" ON example ("url");

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO nobody;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO nobody;
