-- +migrate Up
ALTER TABLE "example" ADD COLUMN "final_url" TEXT DEFAULT '' NOT NULL;
UPDATE "example" SET "final_url" = "url";
ALTER TABLE "example" ALTER COLUMN "final_url" DROP DEFAULT;

ALTER TABLE "example" ADD COLUMN "title" TEXT;
ALTER TABLE "example" ADD COLUMN "description" TEXT;
ALTER TABLE "example" ADD COLUMN "og_description" TEXT;
ALTER TABLE "example" ADD COLUMN "og_type" TEXT;
ALTER TABLE "example" ADD COLUMN "og_image" TEXT;
ALTER TABLE "example" ADD COLUMN "body" TEXT;
ALTER TABLE "example" ADD COLUMN "score" DOUBLE PRECISION DEFAULT 0.0 NOT NULL;
ALTER TABLE "example" ADD COLUMN "is_new" BOOLEAN DEFAULT FALSE NOT NULL;
ALTER TABLE "example" ADD COLUMN "status_code" INT DEFAULT 0 NOT NULL;
ALTER TABLE "example" ADD COLUMN "favicon" TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS "final_url_idx_example" ON example ("final_url");

-- +migrate Down
DROP INDEX "final_url_idx_example";

ALTER TABLE "example" DROP COLUMN "final_url";
ALTER TABLE "example" DROP COLUMN "title";
ALTER TABLE "example" DROP COLUMN "description";
ALTER TABLE "example" DROP COLUMN "og_description";
ALTER TABLE "example" DROP COLUMN "og_type";
ALTER TABLE "example" DROP COLUMN "og_image";
ALTER TABLE "example" DROP COLUMN "body";
ALTER TABLE "example" DROP COLUMN "score";
ALTER TABLE "example" DROP COLUMN "is_new";
ALTER TABLE "example" DROP COLUMN "status_code";
ALTER TABLE "example" DROP COLUMN "favicon";
