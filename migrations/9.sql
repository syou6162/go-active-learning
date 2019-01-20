-- +migrate Up
CREATE TABLE IF NOT EXISTS model (
  "model" TEXT NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS "created_at_model" ON model ("created_at");

-- +migrate Down
DROP INDEX "created_at_model";

DROP TABLE model;
