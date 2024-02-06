/*
  Warnings:

  - You are about to drop the column `label` on the `imei_configuration` table. All the data in the column will be lost.
  - You are about to drop the column `tags` on the `imei_configuration` table. All the data in the column will be lost.
  - You are about to drop the column `label` on the `imsi_configuration` table. All the data in the column will be lost.
  - You are about to drop the column `tags` on the `imsi_configuration` table. All the data in the column will be lost.
  - You are about to drop the column `label` on the `license_plate_configuration` table. All the data in the column will be lost.
  - You are about to drop the column `tags` on the `license_plate_configuration` table. All the data in the column will be lost.
  - You are about to drop the column `tags` on the `mobile_device_configuration` table. All the data in the column will be lost.
  - You are about to drop the column `tags` on the `station_location` table. All the data in the column will be lost.
  - You are about to drop the column `label` on the `vehicle_imsi` table. All the data in the column will be lost.
  - You are about to drop the column `label` on the `vehicle_license_plate` table. All the data in the column will be lost.
  - You are about to drop the column `tags` on the `vehicle_license_plate` table. All the data in the column will be lost.
  - Added the required column `department` to the `station_location` table without a default value. This is not possible if the table is not empty.
  - Added the required column `accuracy` to the `vehicle_license_plate` table without a default value. This is not possible if the table is not empty.

*/
-- DropIndex
DROP INDEX "imsi_configuration_tags_idx";

-- DropIndex
DROP INDEX "license_plate_configuration_tags_idx";

-- DropIndex
DROP INDEX "mobile_device_configuration_tags_idx";

-- DropIndex
DROP INDEX "station_location_tags_idx";

-- DropIndex
DROP INDEX "vehicle_license_plate_tags_idx";

-- AlterTable
ALTER TABLE "imei_configuration" DROP COLUMN "label",
DROP COLUMN "tags",
ADD COLUMN     "permittedLabel" "DevicePermittedLabel" NOT NULL DEFAULT 'NONE';

-- AlterTable
ALTER TABLE "imsi_configuration" DROP COLUMN "label",
DROP COLUMN "tags",
ADD COLUMN     "permittedLabel" "DevicePermittedLabel" NOT NULL DEFAULT 'NONE';

-- AlterTable
ALTER TABLE "license_plate_configuration" DROP COLUMN "label",
DROP COLUMN "tags",
ADD COLUMN     "permittedLabel" "DevicePermittedLabel" NOT NULL DEFAULT 'NONE';

-- AlterTable
ALTER TABLE "mobile_device_configuration" DROP COLUMN "tags";

-- AlterTable
ALTER TABLE "station_location" DROP COLUMN "tags",
ADD COLUMN     "department" TEXT NOT NULL,
ADD COLUMN     "remark" TEXT;

-- AlterTable
ALTER TABLE "vehicle_imsi" DROP COLUMN "label",
ADD COLUMN     "permittedLabel" "DevicePermittedLabel" NOT NULL DEFAULT 'NONE';

-- AlterTable
ALTER TABLE "vehicle_license_plate" DROP COLUMN "label",
DROP COLUMN "tags",
ADD COLUMN     "accuracy" DECIMAL(65,30) NOT NULL,
ADD COLUMN     "permittedLabel" "DevicePermittedLabel" DEFAULT 'NONE';

-- CreateTable
CREATE TABLE "tags" (
    "id" UUID NOT NULL,
    "projectId" UUID NOT NULL,
    "title" TEXT NOT NULL,
    "createdBy" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "tags_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "imsi_configuration_tag" (
    "id" UUID NOT NULL,
    "imsiConfigurationId" UUID NOT NULL,
    "tagId" UUID NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "imsi_configuration_tag_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "imei_configuration_tag" (
    "id" UUID NOT NULL,
    "imeiConfigurationId" UUID NOT NULL,
    "tagId" UUID NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "imei_configuration_tag_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "mobile_device_configuration_tag" (
    "id" UUID NOT NULL,
    "mobileDeviceConfigurationId" UUID NOT NULL,
    "tagId" UUID NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "mobile_device_configuration_tag_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "station_location_tag" (
    "id" UUID NOT NULL,
    "stationLocationId" UUID NOT NULL,
    "tagId" UUID NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "station_location_tag_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "license_plate_configuration_tag" (
    "id" UUID NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "licensePlateConfigurationId" UUID NOT NULL,
    "tagId" UUID NOT NULL,

    CONSTRAINT "license_plate_configuration_tag_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "vehicle_license_plate_tag" (
    "id" UUID NOT NULL,
    "vehicleLicensePlateVehicleId" UUID NOT NULL,
    "tagId" UUID NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "vehicle_license_plate_tag_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "vehicle_tag" (
    "id" UUID NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "tagId" UUID NOT NULL,
    "vehicleId" UUID NOT NULL,

    CONSTRAINT "vehicle_tag_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE INDEX "tags_title_idx" ON "tags" USING GIN ("title" gin_trgm_ops);

-- AddForeignKey
ALTER TABLE "tags" ADD CONSTRAINT "tags_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES "project"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "imsi_configuration_tag" ADD CONSTRAINT "imsi_configuration_tag_imsiConfigurationId_fkey" FOREIGN KEY ("imsiConfigurationId") REFERENCES "imsi_configuration"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "imsi_configuration_tag" ADD CONSTRAINT "imsi_configuration_tag_tagId_fkey" FOREIGN KEY ("tagId") REFERENCES "tags"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "imei_configuration_tag" ADD CONSTRAINT "imei_configuration_tag_imeiConfigurationId_fkey" FOREIGN KEY ("imeiConfigurationId") REFERENCES "imei_configuration"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "imei_configuration_tag" ADD CONSTRAINT "imei_configuration_tag_tagId_fkey" FOREIGN KEY ("tagId") REFERENCES "tags"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "mobile_device_configuration_tag" ADD CONSTRAINT "mobile_device_configuration_tag_mobileDeviceConfigurationI_fkey" FOREIGN KEY ("mobileDeviceConfigurationId") REFERENCES "mobile_device_configuration"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "mobile_device_configuration_tag" ADD CONSTRAINT "mobile_device_configuration_tag_tagId_fkey" FOREIGN KEY ("tagId") REFERENCES "tags"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "station_location_tag" ADD CONSTRAINT "station_location_tag_stationLocationId_fkey" FOREIGN KEY ("stationLocationId") REFERENCES "station_location"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "station_location_tag" ADD CONSTRAINT "station_location_tag_tagId_fkey" FOREIGN KEY ("tagId") REFERENCES "tags"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "license_plate_configuration_tag" ADD CONSTRAINT "license_plate_configuration_tag_licensePlateConfigurationI_fkey" FOREIGN KEY ("licensePlateConfigurationId") REFERENCES "license_plate_configuration"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "license_plate_configuration_tag" ADD CONSTRAINT "license_plate_configuration_tag_tagId_fkey" FOREIGN KEY ("tagId") REFERENCES "tags"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle_license_plate_tag" ADD CONSTRAINT "vehicle_license_plate_tag_vehicleLicensePlateVehicleId_fkey" FOREIGN KEY ("vehicleLicensePlateVehicleId") REFERENCES "vehicle_license_plate"("vehicleId") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle_license_plate_tag" ADD CONSTRAINT "vehicle_license_plate_tag_tagId_fkey" FOREIGN KEY ("tagId") REFERENCES "tags"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle_tag" ADD CONSTRAINT "vehicle_tag_tagId_fkey" FOREIGN KEY ("tagId") REFERENCES "tags"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle_tag" ADD CONSTRAINT "vehicle_tag_vehicleId_fkey" FOREIGN KEY ("vehicleId") REFERENCES "vehicle"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
