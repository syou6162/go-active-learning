-- +migrate Up
ALTER TABLE "tweet" ADD COLUMN "score" DOUBLE PRECISION DEFAULT 0.0 NOT NULL;

-- +migrate Down
ALTER TABLE "tweet" DROP COLUMN "score";
