/*
  Warnings:

  - You are about to drop the `station_health_check_activity` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropForeignKey
ALTER TABLE "station_health_check_activity" DROP CONSTRAINT "station_health_check_activity_stationId_fkey";

-- DropForeignKey
ALTER TABLE "station_location" DROP CONSTRAINT "station_location_currentHealthCheckId_fkey";

-- DropTable
DROP TABLE "station_health_check_activity";

-- CreateTable
CREATE TABLE "station_location_health_check_activity" (
    "id" UUID NOT NULL,
    "stationId" UUID NOT NULL,
    "stationStatus" "StationStatus" NOT NULL,
    "startDatetime" TIMESTAMPTZ(3) NOT NULL,
    "endDatetime" TIMESTAMPTZ(3),
    "createdBy" TEXT NOT NULL,
    "createdAt" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedBy" TEXT,
    "updatedAt" TIMESTAMPTZ(3),

    CONSTRAINT "station_location_health_check_activity_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "station_location" ADD CONSTRAINT "station_location_currentHealthCheckId_fkey" FOREIGN KEY ("currentHealthCheckId") REFERENCES "station_location_health_check_activity"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "station_location_health_check_activity" ADD CONSTRAINT "station_location_health_check_activity_stationId_fkey" FOREIGN KEY ("stationId") REFERENCES "station_location"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
