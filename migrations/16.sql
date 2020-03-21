-- +migrate Up
CREATE TABLE IF NOT EXISTS top_accessed_example (
  "example_id" SERIAL NOT NULL,
  CONSTRAINT top_accessed_example_example_id_fkey FOREIGN KEY ("example_id") REFERENCES example("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS "example_id_idx_top_accessed_example" ON top_accessed_example ("example_id");

-- +migrate Down
DROP INDEX "example_id_idx_top_accessed_example";

DROP TABLE top_accessed_example;
