/*
  Warnings:

  - The `label` column on the `imei_configuration` table would be dropped and recreated. This will lead to data loss if there is data in the column.
  - The `label` column on the `imsi_configuration` table would be dropped and recreated. This will lead to data loss if there is data in the column.
  - The `label` column on the `license_plate_configuration` table would be dropped and recreated. This will lead to data loss if there is data in the column.
  - The `label` column on the `mobile_device_configuration` table would be dropped and recreated. This will lead to data loss if there is data in the column.
  - You are about to drop the column `inspectionLocationId` on the `vehicle` table. All the data in the column will be lost.
  - You are about to drop the column `label` on the `vehicle_image` table. All the data in the column will be lost.
  - The `label` column on the `vehicle_imsi` table would be dropped and recreated. This will lead to data loss if there is data in the column.
  - The `label` column on the `vehicle_license_plate` table would be dropped and recreated. This will lead to data loss if there is data in the column.
  - You are about to drop the `inspection_location` table. If the table is not empty, all the data it contains will be lost.
  - Added the required column `stationLocationId` to the `imei_configuration` table without a default value. This is not possible if the table is not empty.
  - Added the required column `stationLocationId` to the `imsi_configuration` table without a default value. This is not possible if the table is not empty.
  - Added the required column `stationLocationId` to the `vehicle` table without a default value. This is not possible if the table is not empty.

*/
-- CreateEnum
CREATE TYPE "BlacklistPriority" AS ENUM ('WARNING', 'DANGER', 'NORMAL');

-- CreateEnum
CREATE TYPE "DevicePermittedLabel" AS ENUM ('WHITELIST', 'BLACKLIST', 'NONE');

-- DropForeignKey
ALTER TABLE "inspection_location" DROP CONSTRAINT "inspection_location_projectId_fkey";

-- DropForeignKey
ALTER TABLE "vehicle" DROP CONSTRAINT "vehicle_inspectionLocationId_fkey";

-- AlterTable
ALTER TABLE "imei_configuration" ADD COLUMN     "priority" "BlacklistPriority" NOT NULL DEFAULT 'NORMAL',
ADD COLUMN     "stationLocationId" UUID NOT NULL,
DROP COLUMN "label",
ADD COLUMN     "label" "DevicePermittedLabel" NOT NULL DEFAULT 'NONE';

-- AlterTable
ALTER TABLE "imsi_configuration" ADD COLUMN     "priority" "BlacklistPriority" NOT NULL DEFAULT 'NORMAL',
ADD COLUMN     "stationLocationId" UUID NOT NULL,
DROP COLUMN "label",
ADD COLUMN     "label" "DevicePermittedLabel" NOT NULL DEFAULT 'NONE';

-- AlterTable
ALTER TABLE "license_plate_configuration" DROP COLUMN "label",
ADD COLUMN     "label" "DevicePermittedLabel" NOT NULL DEFAULT 'NONE';

-- AlterTable
ALTER TABLE "mobile_device_configuration" DROP COLUMN "label",
ADD COLUMN     "label" "DevicePermittedLabel" NOT NULL DEFAULT 'NONE';

-- AlterTable
ALTER TABLE "vehicle" DROP COLUMN "inspectionLocationId",
ADD COLUMN     "stationLocationId" UUID NOT NULL;

-- AlterTable
ALTER TABLE "vehicle_image" DROP COLUMN "label";

-- AlterTable
ALTER TABLE "vehicle_imsi" DROP COLUMN "label",
ADD COLUMN     "label" "DevicePermittedLabel" NOT NULL DEFAULT 'NONE';

-- AlterTable
ALTER TABLE "vehicle_license_plate" DROP COLUMN "label",
ADD COLUMN     "label" "DevicePermittedLabel" NOT NULL DEFAULT 'NONE';

-- DropTable
DROP TABLE "inspection_location";

-- DropEnum
DROP TYPE "DeviceCaptureLabel";

-- CreateTable
CREATE TABLE "station_location" (
    "id" UUID NOT NULL,
    "projectId" UUID NOT NULL,
    "title" TEXT NOT NULL,
    "description" TEXT,
    "tags" TEXT[] DEFAULT ARRAY[]::TEXT[],
    "latitude" DECIMAL(8,6) NOT NULL,
    "longtitude" DECIMAL(9,6) NOT NULL,
    "createdBy" TEXT NOT NULL,
    "updatedBy" TEXT,
    "deletedBy" TEXT,
    "deletedAt" TIMESTAMP(3),
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3),

    CONSTRAINT "station_location_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE INDEX "station_location_tags_idx" ON "station_location" USING GIN ("tags");

-- CreateIndex
CREATE INDEX "station_location_title_idx" ON "station_location" USING GIN ("title" gin_trgm_ops);

-- AddForeignKey
ALTER TABLE "imsi_configuration" ADD CONSTRAINT "imsi_configuration_stationLocationId_fkey" FOREIGN KEY ("stationLocationId") REFERENCES "station_location"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "imei_configuration" ADD CONSTRAINT "imei_configuration_stationLocationId_fkey" FOREIGN KEY ("stationLocationId") REFERENCES "station_location"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "station_location" ADD CONSTRAINT "station_location_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES "project"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle" ADD CONSTRAINT "vehicle_stationLocationId_fkey" FOREIGN KEY ("stationLocationId") REFERENCES "station_location"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
