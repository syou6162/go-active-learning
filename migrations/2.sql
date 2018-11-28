-- +migrate Up
CREATE TABLE IF NOT EXISTS feature (
  "example_id" SERIAL NOT NULL,
  "feature" TEXT NOT NULL,
  CONSTRAINT feature_example_id_fkey FOREIGN KEY ("example_id") REFERENCES example("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS "example_id_idx_example" ON feature ("example_id");

-- +migrate Down
DROP INDEX "example_id_idx_example";
DROP TABLE feature;
