-- +migrate Up
ALTER TABLE "tweet" ADD COLUMN "label" INT NOT NULL DEFAULT 0;

-- +migrate Down
ALTER TABLE "tweet" DROP COLUMN "label";
