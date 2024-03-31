/*
  Warnings:

  - You are about to drop the `StationLocationConfiguration` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropForeignKey
ALTER TABLE "StationLocationConfiguration" DROP CONSTRAINT "StationLocationConfiguration_stationLocationId_fkey";

-- DropTable
DROP TABLE "StationLocationConfiguration";

-- CreateTable
CREATE TABLE "station_location_configuration" (
    "stationLocationId" UUID NOT NULL,
    "apiKey" TEXT NOT NULL,

    CONSTRAINT "station_location_configuration_pkey" PRIMARY KEY ("stationLocationId")
);

-- CreateIndex
CREATE INDEX "station_location_configuration_apiKey_idx" ON "station_location_configuration"("apiKey");

-- AddForeignKey
ALTER TABLE "station_location_configuration" ADD CONSTRAINT "station_location_configuration_stationLocationId_fkey" FOREIGN KEY ("stationLocationId") REFERENCES "station_location"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
