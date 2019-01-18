-- +migrate Up
ALTER TABLE "tweet" ADD COLUMN "label" INT NOT NULL;

-- +migrate Down
ALTER TABLE "tweet" DROP COLUMN "label";
