/*
  Warnings:

  - Added the required column `stationLocationId` to the `mobile_device_configuration` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "mobile_device_configuration" ADD COLUMN     "stationLocationId" UUID NOT NULL;

-- AddForeignKey
ALTER TABLE "mobile_device_configuration" ADD CONSTRAINT "mobile_device_configuration_stationLocationId_fkey" FOREIGN KEY ("stationLocationId") REFERENCES "station_location"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
