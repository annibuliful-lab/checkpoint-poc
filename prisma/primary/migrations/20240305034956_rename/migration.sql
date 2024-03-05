/*
  Warnings:

  - You are about to drop the `station_device_health_check_activity` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropForeignKey
ALTER TABLE "station_device_health_check_activity" DROP CONSTRAINT "station_device_health_check_activity_stationDeviceId_fkey";

-- DropTable
DROP TABLE "station_device_health_check_activity";

-- CreateTable
CREATE TABLE "station_location_device_health_check_activity" (
    "id" UUID NOT NULL,
    "stationDeviceId" UUID NOT NULL,
    "status" "DeviceStatus" NOT NULL,
    "activityTime" TIMESTAMPTZ(3) NOT NULL,
    "issue" TEXT,
    "createdAt" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "station_location_device_health_check_activity_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "station_location_device_health_check_activity" ADD CONSTRAINT "station_location_device_health_check_activity_stationDevic_fkey" FOREIGN KEY ("stationDeviceId") REFERENCES "station_device"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
