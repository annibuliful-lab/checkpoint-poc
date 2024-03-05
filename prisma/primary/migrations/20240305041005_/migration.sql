-- AlterTable
ALTER TABLE "station_location_device_health_check_activity" ADD COLUMN     "updatedAt" TIMESTAMPTZ(3),
ADD COLUMN     "updatedBy" TEXT;
