-- +migrate Up
ALTER TABLE "model" ADD COLUMN "model_type" INT NOT NULL DEFAULT 0;
ALTER TABLE "model" ADD COLUMN "c" DOUBLE PRECISION DEFAULT 0.0 NOT NULL;
ALTER TABLE "model" ADD COLUMN "accuracy" DOUBLE PRECISION DEFAULT 0.0 NOT NULL;
ALTER TABLE "model" ADD COLUMN "precision" DOUBLE PRECISION DEFAULT 0.0 NOT NULL;
ALTER TABLE "model" ADD COLUMN "recall" DOUBLE PRECISION DEFAULT 0.0 NOT NULL;
ALTER TABLE "model" ADD COLUMN "fvalue" DOUBLE PRECISION DEFAULT 0.0 NOT NULL;

DROP INDEX "created_at_model";
CREATE INDEX IF NOT EXISTS "model_type_created_at_model" ON model ("model_type", "created_at");

-- +migrate Down
ALTER TABLE "model" DROP COLUMN "model_type";
ALTER TABLE "model" DROP COLUMN "c";
ALTER TABLE "model" DROP COLUMN "accuracy";
ALTER TABLE "model" DROP COLUMN "precision";
ALTER TABLE "model" DROP COLUMN "recall";
ALTER TABLE "model" DROP COLUMN "fvalue";

CREATE INDEX IF NOT EXISTS "created_at_model" ON model ("created_at");
DROP INDEX "model_type_created_at_model";
