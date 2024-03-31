/*
  Warnings:

  - Made the column `status` on table `station_vehicle_activity` required. This step will fail if there are existing NULL values in that column.

*/
-- AlterTable
ALTER TABLE "station_vehicle_activity" ALTER COLUMN "status" SET NOT NULL,
ALTER COLUMN "status" SET DEFAULT 'IN_QUEUE';
