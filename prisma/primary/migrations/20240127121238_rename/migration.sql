/*
  Warnings:

  - You are about to drop the `imei_catcher` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `imsi_catcher` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `license_plate` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropForeignKey
ALTER TABLE "imei_catcher" DROP CONSTRAINT "imei_catcher_projectId_fkey";

-- DropForeignKey
ALTER TABLE "imsi_catcher" DROP CONSTRAINT "imsi_catcher_projectId_fkey";

-- DropForeignKey
ALTER TABLE "license_plate" DROP CONSTRAINT "license_plate_projectId_fkey";

-- DropTable
DROP TABLE "imei_catcher";

-- DropTable
DROP TABLE "imsi_catcher";

-- DropTable
DROP TABLE "license_plate";

-- CreateTable
CREATE TABLE "imsi_configuration" (
    "id" UUID NOT NULL,
    "projectId" UUID NOT NULL,
    "imsi" VARCHAR(24) NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3),
    "status" "DeviceCaptureStatus" NOT NULL,

    CONSTRAINT "imsi_configuration_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "imei_configuration" (
    "id" UUID NOT NULL,
    "projectId" UUID NOT NULL,
    "imei" VARCHAR(24) NOT NULL,
    "status" "DeviceCaptureStatus" NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3),

    CONSTRAINT "imei_configuration_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "license_plate_configuration" (
    "id" UUID NOT NULL,
    "projectId" UUID NOT NULL,
    "prefix" TEXT NOT NULL,
    "number" VARCHAR(24) NOT NULL,
    "province" TEXT NOT NULL,
    "country" TEXT,
    "status" "DeviceCaptureStatus" NOT NULL,

    CONSTRAINT "license_plate_configuration_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "imsi_configuration_imsi_key" ON "imsi_configuration"("imsi");

-- CreateIndex
CREATE INDEX "imei_configuration_imei_idx" ON "imei_configuration"("imei");

-- CreateIndex
CREATE INDEX "license_plate_configuration_number_prefix_idx" ON "license_plate_configuration"("number", "prefix");

-- CreateIndex
CREATE INDEX "license_plate_configuration_number_prefix_province_idx" ON "license_plate_configuration"("number", "prefix", "province");

-- AddForeignKey
ALTER TABLE "imsi_configuration" ADD CONSTRAINT "imsi_configuration_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES "project"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "imei_configuration" ADD CONSTRAINT "imei_configuration_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES "project"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "license_plate_configuration" ADD CONSTRAINT "license_plate_configuration_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES "project"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
