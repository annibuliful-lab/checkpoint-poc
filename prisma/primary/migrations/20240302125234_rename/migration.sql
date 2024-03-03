/*
  Warnings:

  - You are about to drop the `station_device_health_check` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropForeignKey
ALTER TABLE "station_device_health_check" DROP CONSTRAINT "station_device_health_check_stationDeviceId_fkey";

-- DropTable
DROP TABLE "station_device_health_check";

-- CreateTable
CREATE TABLE "station_device_health_check_activity" (
    "id" UUID NOT NULL,
    "stationDeviceId" UUID NOT NULL,
    "status" "DeviceStatus" NOT NULL,
    "issue" TEXT,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "station_device_health_check_activity_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "station_device_health_check_activity" ADD CONSTRAINT "station_device_health_check_activity_stationDeviceId_fkey" FOREIGN KEY ("stationDeviceId") REFERENCES "station_device"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
