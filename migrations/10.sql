-- +migrate Up
CREATE TABLE IF NOT EXISTS recommendation (
  "list_type" INT NOT NULL,
  "example_id" SERIAL NOT NULL,
   CONSTRAINT recommendation_example_id_fkey FOREIGN KEY ("example_id") REFERENCES example("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS "list_type_idx_recommendation" ON recommendation ("list_type");

-- +migrate Down
DROP INDEX "list_type_idx_recommendation";

DROP TABLE recommendation;
