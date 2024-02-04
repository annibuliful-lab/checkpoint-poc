/*
  Warnings:

  - Added the required column `type` to the `license_plate_configuration` table without a default value. This is not possible if the table is not empty.
  - Added the required column `createdBy` to the `station_device` table without a default value. This is not possible if the table is not empty.
  - Added the required column `type` to the `vehicle_image` table without a default value. This is not possible if the table is not empty.
  - Added the required column `type` to the `vehicle_license_plate` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "license_plate_configuration" ADD COLUMN     "type" TEXT NOT NULL;

-- AlterTable
ALTER TABLE "station_device" ADD COLUMN     "createdBy" TEXT NOT NULL,
ADD COLUMN     "deletedAt" TIMESTAMP(3),
ADD COLUMN     "deletedBy" TEXT,
ADD COLUMN     "hardwareVersion" TEXT,
ADD COLUMN     "updatedAt" TIMESTAMP(3),
ADD COLUMN     "updatedBy" TEXT;

-- AlterTable
ALTER TABLE "vehicle_image" ADD COLUMN     "type" TEXT NOT NULL;

-- AlterTable
ALTER TABLE "vehicle_license_plate" ADD COLUMN     "type" TEXT NOT NULL;

-- CreateTable
CREATE TABLE "station_officer" (
    "id" UUID NOT NULL,
    "stationLocationId" UUID NOT NULL,
    "projectAccountId" UUID,
    "firstname" TEXT NOT NULL,
    "lastname" TEXT NOT NULL,
    "msisdn" VARCHAR(10) NOT NULL,
    "createdBy" TEXT NOT NULL,
    "updatedBy" TEXT,
    "deletedBy" TEXT,
    "deletedAt" TIMESTAMP(3),
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3),

    CONSTRAINT "station_officer_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE INDEX "station_officer_firstname_lastname_idx" ON "station_officer"("firstname", "lastname");

-- AddForeignKey
ALTER TABLE "station_officer" ADD CONSTRAINT "station_officer_stationLocationId_fkey" FOREIGN KEY ("stationLocationId") REFERENCES "station_location"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "station_officer" ADD CONSTRAINT "station_officer_projectAccountId_fkey" FOREIGN KEY ("projectAccountId") REFERENCES "project_account"("id") ON DELETE SET NULL ON UPDATE CASCADE;
