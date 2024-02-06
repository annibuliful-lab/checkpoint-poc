/*
  Warnings:

  - Changed the type of `type` on the `vehicle_image` table. No cast exists, the column would be dropped and recreated, which cannot be done if there is data, since the column is required.

*/
-- CreateEnum
CREATE TYPE "ImageType" AS ENUM ('FRONT', 'REAR', 'DRIVER', 'LICENSE_PLATE');

-- AlterTable
ALTER TABLE "vehicle_image" DROP COLUMN "type",
ADD COLUMN     "type" "ImageType" NOT NULL;

-- AlterTable
ALTER TABLE "vehicle_license_plate" ADD COLUMN     "s3Key" TEXT,
ALTER COLUMN "label" DROP NOT NULL,
ALTER COLUMN "priority" DROP NOT NULL;

-- CreateIndex
CREATE INDEX "imsi_configuration_tags_idx" ON "imsi_configuration" USING GIN ("tags");

-- CreateIndex
CREATE INDEX "license_plate_configuration_tags_idx" ON "license_plate_configuration" USING GIN ("tags");

-- CreateIndex
CREATE INDEX "vehicle_license_plate_tags_idx" ON "vehicle_license_plate" USING GIN ("tags");
