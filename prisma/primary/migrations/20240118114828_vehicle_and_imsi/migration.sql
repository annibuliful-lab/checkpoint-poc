-- CreateTable
CREATE TABLE "license_plate" (
    "id" UUID NOT NULL,
    "projectId" UUID NOT NULL,
    "licensePlate" VARCHAR(24) NOT NULL,
    "licensePlateProvince" TEXT NOT NULL,
    "isBlacklist" BOOLEAN NOT NULL DEFAULT false,

    CONSTRAINT "license_plate_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "vehicle" (
    "id" UUID NOT NULL,
    "projectId" UUID NOT NULL,
    "inspectionLocationId" UUID NOT NULL,
    "licensePlateId" UUID NOT NULL,
    "brand" TEXT NOT NULL,
    "color" TEXT NOT NULL,
    "country" TEXT,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3),

    CONSTRAINT "vehicle_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "vehicle_image" (
    "id" UUID NOT NULL,
    "vehicleId" UUID NOT NULL,
    "s3Key" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "vehicle_image_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "vehivle_imsi" (
    "id" UUID NOT NULL,
    "imsiCatcherId" UUID NOT NULL,
    "vehicleId" UUID NOT NULL,
    "imei" VARCHAR(24) NOT NULL,
    "signalStrength" DOUBLE PRECISION NOT NULL,
    "registerDateTime" TIMESTAMP(3) NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "vehivle_imsi_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "imsi_catcher" (
    "id" UUID NOT NULL,
    "projectId" UUID NOT NULL,
    "imsi" VARCHAR(24) NOT NULL,
    "isBlacklist" BOOLEAN NOT NULL DEFAULT false,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "imsi_catcher_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE INDEX "license_plate_licensePlate_idx" ON "license_plate"("licensePlate");

-- CreateIndex
CREATE INDEX "license_plate_licensePlateProvince_idx" ON "license_plate"("licensePlateProvince");

-- CreateIndex
CREATE INDEX "vehivle_imsi_imei_idx" ON "vehivle_imsi"("imei");

-- CreateIndex
CREATE INDEX "imsi_catcher_imsi_idx" ON "imsi_catcher"("imsi");

-- AddForeignKey
ALTER TABLE "license_plate" ADD CONSTRAINT "license_plate_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES "project"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle" ADD CONSTRAINT "vehicle_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES "project"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle" ADD CONSTRAINT "vehicle_inspectionLocationId_fkey" FOREIGN KEY ("inspectionLocationId") REFERENCES "inspection_location"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle" ADD CONSTRAINT "vehicle_licensePlateId_fkey" FOREIGN KEY ("licensePlateId") REFERENCES "license_plate"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle_image" ADD CONSTRAINT "vehicle_image_vehicleId_fkey" FOREIGN KEY ("vehicleId") REFERENCES "vehicle"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehivle_imsi" ADD CONSTRAINT "vehivle_imsi_vehicleId_fkey" FOREIGN KEY ("vehicleId") REFERENCES "vehicle"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehivle_imsi" ADD CONSTRAINT "vehivle_imsi_imsiCatcherId_fkey" FOREIGN KEY ("imsiCatcherId") REFERENCES "imsi_catcher"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "imsi_catcher" ADD CONSTRAINT "imsi_catcher_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES "project"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
