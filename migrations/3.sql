-- +migrate Up
ALTER TABLE "example" ALTER COLUMN "title" SET DEFAULT '';
ALTER TABLE "example" ALTER COLUMN "description" SET DEFAULT '';
ALTER TABLE "example" ALTER COLUMN "og_description" SET DEFAULT '';
ALTER TABLE "example" ALTER COLUMN "og_type" SET DEFAULT '';
ALTER TABLE "example" ALTER COLUMN "og_image" SET DEFAULT '';
ALTER TABLE "example" ALTER COLUMN "body" SET DEFAULT '';
ALTER TABLE "example" ALTER COLUMN "favicon" SET DEFAULT '';

-- +migrate Down
ALTER TABLE "example" ALTER COLUMN "title" DROP DEFAULT;
ALTER TABLE "example" ALTER COLUMN "description" DROP DEFAULT;
ALTER TABLE "example" ALTER COLUMN "og_description" DROP DEFAULT;
ALTER TABLE "example" ALTER COLUMN "og_type" DROP DEFAULT;
ALTER TABLE "example" ALTER COLUMN "og_image" DROP DEFAULT;
ALTER TABLE "example" ALTER COLUMN "body" DROP DEFAULT;
ALTER TABLE "example" ALTER COLUMN "favicon" DROP DEFAULT;
