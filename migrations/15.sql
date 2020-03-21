-- +migrate Up
CREATE TABLE IF NOT EXISTS related_example (
  "example_id" SERIAL NOT NULL,
  "related_example_id" SERIAL NOT NULL,
  CONSTRAINT related_example_example_id_fkey FOREIGN KEY ("example_id") REFERENCES example("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT related_example_related_example_id_fkey FOREIGN KEY ("related_example_id") REFERENCES example("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CHECK(example_id != related_example_id)
);

CREATE INDEX IF NOT EXISTS "example_id_idx_related_example" ON related_example ("example_id");

-- +migrate Down
DROP INDEX "example_id_idx_related_example";

DROP TABLE related_example;
