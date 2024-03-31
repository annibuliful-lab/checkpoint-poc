/*
  Warnings:

  - You are about to alter the column `accuracy` on the `vehicle_license_plate` table. The data in that column could be lost. The data in that column will be cast from `Decimal(65,30)` to `DoublePrecision`.
  - Added the required column `reportIssueBy` to the `station_vehicle_activity` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "station_vehicle_activity" ADD COLUMN     "issue" TEXT,
ADD COLUMN     "reportIssueBy" TEXT NOT NULL;

-- AlterTable
ALTER TABLE "vehicle_license_plate" ALTER COLUMN "accuracy" SET DATA TYPE DOUBLE PRECISION;
