/*
  Warnings:

  - You are about to drop the column `vehicleId` on the `vehicle_image` table. All the data in the column will be lost.
  - You are about to drop the column `vehicleId` on the `vehicle_imsi` table. All the data in the column will be lost.
  - The primary key for the `vehicle_license_plate` table will be changed. If it partially fails, the table could be left without primary key constraint.
  - You are about to drop the column `vehicleId` on the `vehicle_license_plate` table. All the data in the column will be lost.
  - You are about to drop the `vehicle` table. If the table is not empty, all the data it contains will be lost.
  - Added the required column `stationVehicleActivityId` to the `vehicle_image` table without a default value. This is not possible if the table is not empty.
  - Added the required column `stationVehicleActivityId` to the `vehicle_imsi` table without a default value. This is not possible if the table is not empty.
  - Added the required column `stationVehicleActivityId` to the `vehicle_license_plate` table without a default value. This is not possible if the table is not empty.

*/
-- DropForeignKey
ALTER TABLE "vehicle" DROP CONSTRAINT "vehicle_projectId_fkey";

-- DropForeignKey
ALTER TABLE "vehicle" DROP CONSTRAINT "vehicle_stationLocationId_fkey";

-- DropForeignKey
ALTER TABLE "vehicle_image" DROP CONSTRAINT "vehicle_image_vehicleId_fkey";

-- DropForeignKey
ALTER TABLE "vehicle_imsi" DROP CONSTRAINT "vehicle_imsi_vehicleId_fkey";

-- DropForeignKey
ALTER TABLE "vehicle_license_plate" DROP CONSTRAINT "vehicle_license_plate_vehicleId_fkey";

-- DropForeignKey
ALTER TABLE "vehicle_license_plate_tag" DROP CONSTRAINT "vehicle_license_plate_tag_vehicleLicensePlateVehicleId_fkey";

-- DropForeignKey
ALTER TABLE "vehicle_tag" DROP CONSTRAINT "vehicle_tag_vehicleId_fkey";

-- AlterTable
ALTER TABLE "vehicle_image" DROP COLUMN "vehicleId",
ADD COLUMN     "stationVehicleActivityId" UUID NOT NULL;

-- AlterTable
ALTER TABLE "vehicle_imsi" DROP COLUMN "vehicleId",
ADD COLUMN     "stationVehicleActivityId" UUID NOT NULL;

-- AlterTable
ALTER TABLE "vehicle_license_plate" DROP CONSTRAINT "vehicle_license_plate_pkey",
DROP COLUMN "vehicleId",
ADD COLUMN     "stationVehicleActivityId" UUID NOT NULL,
ADD CONSTRAINT "vehicle_license_plate_pkey" PRIMARY KEY ("stationVehicleActivityId");

-- DropTable
DROP TABLE "vehicle";

-- CreateTable
CREATE TABLE "station_vehicle_activity" (
    "id" UUID NOT NULL,
    "projectId" UUID NOT NULL,
    "stationLocationId" UUID NOT NULL,
    "brand" TEXT NOT NULL,
    "color" TEXT NOT NULL,
    "model" TEXT NOT NULL,
    "status" "RemarkState",
    "createdBy" TEXT NOT NULL,
    "updatedBy" TEXT,
    "createdAt" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMPTZ(3),

    CONSTRAINT "station_vehicle_activity_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE INDEX "station_vehicle_activity_brand_idx" ON "station_vehicle_activity"("brand");

-- CreateIndex
CREATE INDEX "station_vehicle_activity_color_idx" ON "station_vehicle_activity"("color");

-- CreateIndex
CREATE INDEX "station_vehicle_activity_model_idx" ON "station_vehicle_activity"("model");

-- AddForeignKey
ALTER TABLE "vehicle_license_plate_tag" ADD CONSTRAINT "vehicle_license_plate_tag_vehicleLicensePlateVehicleId_fkey" FOREIGN KEY ("vehicleLicensePlateVehicleId") REFERENCES "vehicle_license_plate"("stationVehicleActivityId") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle_tag" ADD CONSTRAINT "vehicle_tag_vehicleId_fkey" FOREIGN KEY ("vehicleId") REFERENCES "station_vehicle_activity"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "station_vehicle_activity" ADD CONSTRAINT "station_vehicle_activity_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES "project"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "station_vehicle_activity" ADD CONSTRAINT "station_vehicle_activity_stationLocationId_fkey" FOREIGN KEY ("stationLocationId") REFERENCES "station_location"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle_license_plate" ADD CONSTRAINT "vehicle_license_plate_stationVehicleActivityId_fkey" FOREIGN KEY ("stationVehicleActivityId") REFERENCES "station_vehicle_activity"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle_image" ADD CONSTRAINT "vehicle_image_stationVehicleActivityId_fkey" FOREIGN KEY ("stationVehicleActivityId") REFERENCES "station_vehicle_activity"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle_imsi" ADD CONSTRAINT "vehicle_imsi_stationVehicleActivityId_fkey" FOREIGN KEY ("stationVehicleActivityId") REFERENCES "station_vehicle_activity"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
