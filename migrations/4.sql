-- +migrate Up
DROP INDEX "final_url_idx_example";

-- +migrate Down
CREATE UNIQUE INDEX IF NOT EXISTS "final_url_idx_example" ON example ("final_url");
