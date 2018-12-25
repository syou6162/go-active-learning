-- +migrate Up
ALTER TABLE "example" ADD COLUMN "error_count" INT NOT NULL DEFAULT 0;

-- +migrate Down
ALTER TABLE "example" DROP COLUMN "error_count";
