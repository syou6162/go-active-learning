-- +migrate Up
CREATE TABLE IF NOT EXISTS recommendation (
  "list_type" INT NOT NULL,
  "example_id" SERIAL NOT NULL,
   CONSTRAINT recommendation_example_id_fkey FOREIGN KEY ("example_id") REFERENCES example("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS "example_id_idx_recommendation" ON recommendation ("example_id");

-- +migrate Down
DROP INDEX "created_at_model";

DROP TABLE recommendation;
