/*
  Warnings:

  - You are about to drop the `tags` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropForeignKey
ALTER TABLE "imei_configuration_tag" DROP CONSTRAINT "imei_configuration_tag_tagId_fkey";

-- DropForeignKey
ALTER TABLE "imsi_configuration_tag" DROP CONSTRAINT "imsi_configuration_tag_tagId_fkey";

-- DropForeignKey
ALTER TABLE "license_plate_configuration_tag" DROP CONSTRAINT "license_plate_configuration_tag_tagId_fkey";

-- DropForeignKey
ALTER TABLE "mobile_device_configuration_tag" DROP CONSTRAINT "mobile_device_configuration_tag_tagId_fkey";

-- DropForeignKey
ALTER TABLE "station_location_tag" DROP CONSTRAINT "station_location_tag_tagId_fkey";

-- DropForeignKey
ALTER TABLE "tags" DROP CONSTRAINT "tags_projectId_fkey";

-- DropForeignKey
ALTER TABLE "vehicle_license_plate_tag" DROP CONSTRAINT "vehicle_license_plate_tag_tagId_fkey";

-- DropForeignKey
ALTER TABLE "vehicle_tag" DROP CONSTRAINT "vehicle_tag_tagId_fkey";

-- DropTable
DROP TABLE "tags";

-- CreateTable
CREATE TABLE "tag" (
    "id" UUID NOT NULL,
    "projectId" UUID NOT NULL,
    "title" TEXT NOT NULL,
    "createdBy" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "tag_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE INDEX "tag_title_idx" ON "tag" USING GIN ("title" gin_trgm_ops);

-- AddForeignKey
ALTER TABLE "tag" ADD CONSTRAINT "tag_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES "project"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "imsi_configuration_tag" ADD CONSTRAINT "imsi_configuration_tag_tagId_fkey" FOREIGN KEY ("tagId") REFERENCES "tag"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "imei_configuration_tag" ADD CONSTRAINT "imei_configuration_tag_tagId_fkey" FOREIGN KEY ("tagId") REFERENCES "tag"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "mobile_device_configuration_tag" ADD CONSTRAINT "mobile_device_configuration_tag_tagId_fkey" FOREIGN KEY ("tagId") REFERENCES "tag"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "station_location_tag" ADD CONSTRAINT "station_location_tag_tagId_fkey" FOREIGN KEY ("tagId") REFERENCES "tag"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "license_plate_configuration_tag" ADD CONSTRAINT "license_plate_configuration_tag_tagId_fkey" FOREIGN KEY ("tagId") REFERENCES "tag"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle_license_plate_tag" ADD CONSTRAINT "vehicle_license_plate_tag_tagId_fkey" FOREIGN KEY ("tagId") REFERENCES "tag"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle_tag" ADD CONSTRAINT "vehicle_tag_tagId_fkey" FOREIGN KEY ("tagId") REFERENCES "tag"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
