-- CreateTable
CREATE TABLE "StationLocationConfiguration" (
    "stationLocationId" UUID NOT NULL,
    "apiKey" TEXT NOT NULL,

    CONSTRAINT "StationLocationConfiguration_pkey" PRIMARY KEY ("stationLocationId")
);

-- AddForeignKey
ALTER TABLE "StationLocationConfiguration" ADD CONSTRAINT "StationLocationConfiguration_stationLocationId_fkey" FOREIGN KEY ("stationLocationId") REFERENCES "station_location"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
