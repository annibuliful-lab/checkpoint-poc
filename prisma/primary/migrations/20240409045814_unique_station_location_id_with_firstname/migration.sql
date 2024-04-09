/*
  Warnings:

  - A unique constraint covering the columns `[stationLocationId,firstname]` on the table `station_officer` will be added. If there are existing duplicate values, this will fail.

*/
-- CreateIndex
CREATE UNIQUE INDEX "station_officer_stationLocationId_firstname_key" ON "station_officer"("stationLocationId", "firstname");
