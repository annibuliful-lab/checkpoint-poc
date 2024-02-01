/*
  Warnings:

  - You are about to drop the column `deletedAt` on the `vehicle` table. All the data in the column will be lost.
  - You are about to drop the column `deletedBy` on the `vehicle` table. All the data in the column will be lost.
  - Added the required column `mcc` to the `imsi_configuration` table without a default value. This is not possible if the table is not empty.
  - Added the required column `mnc` to the `imsi_configuration` table without a default value. This is not possible if the table is not empty.
  - Added the required column `mcc` to the `vehicle_imsi` table without a default value. This is not possible if the table is not empty.
  - Added the required column `mnc` to the `vehicle_imsi` table without a default value. This is not possible if the table is not empty.

*/
-- CreateEnum
CREATE TYPE "DeviceStatus" AS ENUM ('ONLINE', 'OFFLINE');

-- CreateEnum
CREATE TYPE "StationStatus" AS ENUM ('ONLINE', 'OFFLINE', 'CLOSED', 'MAINTENANCE');

-- AlterTable
ALTER TABLE "imei_configuration" ADD COLUMN     "tags" TEXT[];

-- AlterTable
ALTER TABLE "imsi_configuration" ADD COLUMN     "mcc" VARCHAR(3) NOT NULL,
ADD COLUMN     "mnc" VARCHAR(3) NOT NULL,
ADD COLUMN     "tags" TEXT[];

-- AlterTable
ALTER TABLE "license_plate_configuration" ADD COLUMN     "blacklistPriority" "BlacklistPriority" NOT NULL DEFAULT 'NORMAL',
ADD COLUMN     "tags" TEXT[];

-- AlterTable
ALTER TABLE "project_account" ADD COLUMN     "firstname" TEXT,
ADD COLUMN     "lastname" TEXT,
ADD COLUMN     "msisdn" TEXT;

-- AlterTable
ALTER TABLE "station_location" ADD COLUMN     "currentHealthCheckId" UUID,
ADD COLUMN     "stationStatus" "StationStatus",
ALTER COLUMN "tags" DROP DEFAULT;

-- AlterTable
ALTER TABLE "vehicle" DROP COLUMN "deletedAt",
DROP COLUMN "deletedBy";

-- AlterTable
ALTER TABLE "vehicle_imsi" ADD COLUMN     "mcc" VARCHAR(3) NOT NULL,
ADD COLUMN     "mnc" VARCHAR(3) NOT NULL;

-- AlterTable
ALTER TABLE "vehicle_license_plate" ADD COLUMN     "priority" "BlacklistPriority" NOT NULL DEFAULT 'NORMAL',
ADD COLUMN     "tags" TEXT[];

-- CreateTable
CREATE TABLE "station_device" (
    "id" UUID NOT NULL,
    "stationLocationId" UUID NOT NULL,
    "title" TEXT NOT NULL,
    "softwareVersion" TEXT,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "station_device_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "station_device_health_check" (
    "id" UUID NOT NULL,
    "stationDeviceId" UUID NOT NULL,
    "status" "DeviceStatus" NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3),

    CONSTRAINT "station_device_health_check_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "station_health_check_activity" (
    "id" UUID NOT NULL,
    "stationId" UUID NOT NULL,
    "stationStatus" "StationStatus" NOT NULL,
    "startDatetime" TIMESTAMP(3) NOT NULL,
    "endDatetime" TIMESTAMP(3),
    "createdBy" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedBy" TEXT,
    "updatedAt" TIMESTAMP(3),

    CONSTRAINT "station_health_check_activity_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "config_line_notify" (
    "id" UUID NOT NULL,
    "projectId" UUID NOT NULL,
    "accountId" UUID NOT NULL,
    "name" TEXT NOT NULL,
    "note" TEXT,
    "token" TEXT NOT NULL,
    "active" BOOLEAN NOT NULL DEFAULT false,
    "createdBy" TEXT NOT NULL,
    "updatedBy" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "config_line_notify_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "config_line_notify_token_key" ON "config_line_notify"("token");

-- CreateIndex
CREATE INDEX "vehicle_imsi_mnc_mcc_idx" ON "vehicle_imsi"("mnc", "mcc");

-- AddForeignKey
ALTER TABLE "station_location" ADD CONSTRAINT "station_location_currentHealthCheckId_fkey" FOREIGN KEY ("currentHealthCheckId") REFERENCES "station_health_check_activity"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "station_device" ADD CONSTRAINT "station_device_stationLocationId_fkey" FOREIGN KEY ("stationLocationId") REFERENCES "station_location"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "station_device_health_check" ADD CONSTRAINT "station_device_health_check_stationDeviceId_fkey" FOREIGN KEY ("stationDeviceId") REFERENCES "station_device"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "station_health_check_activity" ADD CONSTRAINT "station_health_check_activity_stationId_fkey" FOREIGN KEY ("stationId") REFERENCES "station_location"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "config_line_notify" ADD CONSTRAINT "config_line_notify_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES "project"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
