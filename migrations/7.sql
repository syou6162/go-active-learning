-- +migrate Up
ALTER TABLE "tweet" ADD COLUMN "retweeted" BOOLEAN NOT NULL DEFAULT false;

-- +migrate Down
ALTER TABLE "tweet" DROP COLUMN "retweeted";
