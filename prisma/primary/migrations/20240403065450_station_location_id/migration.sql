/*
  Warnings:

  - Added the required column `stationLocationId` to the `vehicle_target_configuration` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "vehicle_target_configuration" ADD COLUMN     "stationLocationId" UUID NOT NULL;

-- AddForeignKey
ALTER TABLE "vehicle_target_configuration" ADD CONSTRAINT "vehicle_target_configuration_stationLocationId_fkey" FOREIGN KEY ("stationLocationId") REFERENCES "station_location"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
