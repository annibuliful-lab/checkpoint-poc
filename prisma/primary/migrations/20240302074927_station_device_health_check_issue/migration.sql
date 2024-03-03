/*
  Warnings:

  - You are about to drop the column `updatedAt` on the `station_device_health_check` table. All the data in the column will be lost.

*/
-- AlterTable
ALTER TABLE "station_device_health_check" DROP COLUMN "updatedAt",
ADD COLUMN     "issue" TEXT;
