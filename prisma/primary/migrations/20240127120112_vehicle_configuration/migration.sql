/*
  Warnings:

  - You are about to drop the column `isBlacklist` on the `imsi_catcher` table. All the data in the column will be lost.
  - You are about to drop the column `isBlacklist` on the `license_plate` table. All the data in the column will be lost.
  - You are about to drop the column `licensePlate` on the `license_plate` table. All the data in the column will be lost.
  - You are about to drop the column `licensePlateProvince` on the `license_plate` table. All the data in the column will be lost.
  - You are about to drop the column `country` on the `vehicle` table. All the data in the column will be lost.
  - You are about to drop the column `licensePlateId` on the `vehicle` table. All the data in the column will be lost.
  - You are about to drop the `vehivle_imsi` table. If the table is not empty, all the data it contains will be lost.
  - A unique constraint covering the columns `[imsi]` on the table `imsi_catcher` will be added. If there are existing duplicate values, this will fail.
  - Added the required column `status` to the `imsi_catcher` table without a default value. This is not possible if the table is not empty.
  - Added the required column `number` to the `license_plate` table without a default value. This is not possible if the table is not empty.
  - Added the required column `prefix` to the `license_plate` table without a default value. This is not possible if the table is not empty.
  - Added the required column `province` to the `license_plate` table without a default value. This is not possible if the table is not empty.
  - Added the required column `status` to the `license_plate` table without a default value. This is not possible if the table is not empty.

*/
-- CreateEnum
CREATE TYPE "DeviceCaptureStatus" AS ENUM ('WHITELIST', 'BLACKLIST');

-- CreateEnum
CREATE TYPE "RemarkState" AS ENUM ('WHITELIST', 'BLACKLIST', 'IN_QUEUE', 'IN_PROGRESS', 'PASSED', 'WAITING', 'INVESTIGATING', 'SUSPICION');

-- DropForeignKey
ALTER TABLE "vehicle" DROP CONSTRAINT "vehicle_licensePlateId_fkey";

-- DropForeignKey
ALTER TABLE "vehivle_imsi" DROP CONSTRAINT "vehivle_imsi_imsiCatcherId_fkey";

-- DropForeignKey
ALTER TABLE "vehivle_imsi" DROP CONSTRAINT "vehivle_imsi_vehicleId_fkey";

-- DropIndex
DROP INDEX "imsi_catcher_imsi_idx";

-- DropIndex
DROP INDEX "license_plate_licensePlateProvince_idx";

-- DropIndex
DROP INDEX "license_plate_licensePlate_idx";

-- AlterTable
ALTER TABLE "imsi_catcher" DROP COLUMN "isBlacklist",
ADD COLUMN     "status" "DeviceCaptureStatus" NOT NULL,
ADD COLUMN     "updatedAt" TIMESTAMP(3);

-- AlterTable
ALTER TABLE "license_plate" DROP COLUMN "isBlacklist",
DROP COLUMN "licensePlate",
DROP COLUMN "licensePlateProvince",
ADD COLUMN     "country" TEXT,
ADD COLUMN     "number" VARCHAR(24) NOT NULL,
ADD COLUMN     "prefix" TEXT NOT NULL,
ADD COLUMN     "province" TEXT NOT NULL,
ADD COLUMN     "status" "DeviceCaptureStatus" NOT NULL;

-- AlterTable
ALTER TABLE "vehicle" DROP COLUMN "country",
DROP COLUMN "licensePlateId",
ADD COLUMN     "status" "RemarkState" NOT NULL DEFAULT 'IN_PROGRESS';

-- AlterTable
ALTER TABLE "vehicle_image" ADD COLUMN     "status" "DeviceCaptureStatus";

-- DropTable
DROP TABLE "vehivle_imsi";

-- CreateTable
CREATE TABLE "imei_catcher" (
    "id" UUID NOT NULL,
    "projectId" UUID NOT NULL,
    "imei" VARCHAR(24) NOT NULL,
    "status" "DeviceCaptureStatus" NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3),

    CONSTRAINT "imei_catcher_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "vehicle_license_plate" (
    "vehicleId" UUID NOT NULL,
    "prefix" VARCHAR(4) NOT NULL,
    "number" VARCHAR(8) NOT NULL,
    "province" VARCHAR(32) NOT NULL,
    "country" TEXT,
    "status" "DeviceCaptureStatus",

    CONSTRAINT "vehicle_license_plate_pkey" PRIMARY KEY ("vehicleId")
);

-- CreateTable
CREATE TABLE "vehicle_imsi" (
    "id" UUID NOT NULL,
    "vehicleId" UUID NOT NULL,
    "labelStatus" "RemarkState" NOT NULL DEFAULT 'IN_PROGRESS',
    "imsi" VARCHAR(24) NOT NULL,
    "imei" VARCHAR(24) NOT NULL,
    "signalStrength" DOUBLE PRECISION NOT NULL,
    "registerDateTime" TIMESTAMP(3) NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3),

    CONSTRAINT "vehicle_imsi_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE INDEX "imei_catcher_imei_idx" ON "imei_catcher"("imei");

-- CreateIndex
CREATE INDEX "vehicle_imsi_imei_imsi_idx" ON "vehicle_imsi"("imei", "imsi");

-- CreateIndex
CREATE UNIQUE INDEX "imsi_catcher_imsi_key" ON "imsi_catcher"("imsi");

-- CreateIndex
CREATE INDEX "license_plate_number_prefix_idx" ON "license_plate"("number", "prefix");

-- CreateIndex
CREATE INDEX "license_plate_number_prefix_province_idx" ON "license_plate"("number", "prefix", "province");

-- AddForeignKey
ALTER TABLE "imei_catcher" ADD CONSTRAINT "imei_catcher_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES "project"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle_license_plate" ADD CONSTRAINT "vehicle_license_plate_vehicleId_fkey" FOREIGN KEY ("vehicleId") REFERENCES "vehicle"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle_imsi" ADD CONSTRAINT "vehicle_imsi_vehicleId_fkey" FOREIGN KEY ("vehicleId") REFERENCES "vehicle"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
