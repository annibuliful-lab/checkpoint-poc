/*
  Warnings:

  - Added the required column `activityTime` to the `station_device_health_check_activity` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "station_device_health_check_activity" ADD COLUMN     "activityTime" TIMESTAMP(3) NOT NULL;
