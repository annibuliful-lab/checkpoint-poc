/*
  Warnings:

  - Made the column `permittedLabel` on table `vehicle_license_plate` required. This step will fail if there are existing NULL values in that column.
  - Made the column `blacklistPriority` on table `vehicle_license_plate` required. This step will fail if there are existing NULL values in that column.

*/
-- AlterTable
ALTER TABLE "vehicle_license_plate" ALTER COLUMN "permittedLabel" SET NOT NULL,
ALTER COLUMN "blacklistPriority" SET NOT NULL;
