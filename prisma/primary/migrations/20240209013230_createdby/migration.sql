/*
  Warnings:

  - Added the required column `createdBy` to the `imei_configuration_tag` table without a default value. This is not possible if the table is not empty.
  - Added the required column `createdBy` to the `imsi_configuration_tag` table without a default value. This is not possible if the table is not empty.
  - Added the required column `createdBy` to the `license_plate_configuration_tag` table without a default value. This is not possible if the table is not empty.
  - Added the required column `createdBy` to the `mobile_device_configuration_tag` table without a default value. This is not possible if the table is not empty.
  - Added the required column `createdBy` to the `station_location_tag` table without a default value. This is not possible if the table is not empty.
  - Added the required column `createdBy` to the `vehicle_license_plate_tag` table without a default value. This is not possible if the table is not empty.
  - Added the required column `createdBy` to the `vehicle_tag` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "imei_configuration_tag" ADD COLUMN     "createdBy" TEXT NOT NULL;

-- AlterTable
ALTER TABLE "imsi_configuration_tag" ADD COLUMN     "createdBy" TEXT NOT NULL;

-- AlterTable
ALTER TABLE "license_plate_configuration_tag" ADD COLUMN     "createdBy" TEXT NOT NULL;

-- AlterTable
ALTER TABLE "mobile_device_configuration_tag" ADD COLUMN     "createdBy" TEXT NOT NULL;

-- AlterTable
ALTER TABLE "station_location_tag" ADD COLUMN     "createdBy" TEXT NOT NULL;

-- AlterTable
ALTER TABLE "vehicle_license_plate_tag" ADD COLUMN     "createdBy" TEXT NOT NULL;

-- AlterTable
ALTER TABLE "vehicle_tag" ADD COLUMN     "createdBy" TEXT NOT NULL;
