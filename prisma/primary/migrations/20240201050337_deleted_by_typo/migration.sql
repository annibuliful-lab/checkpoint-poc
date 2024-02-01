/*
  Warnings:

  - You are about to drop the column `deleteBy` on the `account` table. All the data in the column will be lost.
  - You are about to drop the column `deleteBy` on the `imei_configuration` table. All the data in the column will be lost.
  - You are about to drop the column `deleteBy` on the `imsi_configuration` table. All the data in the column will be lost.
  - You are about to drop the column `deleteBy` on the `inspection_location` table. All the data in the column will be lost.
  - You are about to drop the column `deleteBy` on the `license_plate_configuration` table. All the data in the column will be lost.
  - You are about to drop the column `deleteBy` on the `mobile_device_configuration` table. All the data in the column will be lost.
  - You are about to drop the column `deleteBy` on the `project` table. All the data in the column will be lost.
  - You are about to drop the column `deleteBy` on the `project_account` table. All the data in the column will be lost.
  - You are about to drop the column `deleteBy` on the `vehicle` table. All the data in the column will be lost.
  - You are about to drop the column `deleteBy` on the `vehicle_image` table. All the data in the column will be lost.
  - You are about to drop the column `deleteBy` on the `vehicle_imsi` table. All the data in the column will be lost.
  - You are about to drop the column `deleteBy` on the `vehicle_license_plate` table. All the data in the column will be lost.

*/
-- AlterTable
ALTER TABLE "account" DROP COLUMN "deleteBy",
ADD COLUMN     "DeletedBy" TEXT;

-- AlterTable
ALTER TABLE "imei_configuration" DROP COLUMN "deleteBy",
ADD COLUMN     "DeletedBy" TEXT;

-- AlterTable
ALTER TABLE "imsi_configuration" DROP COLUMN "deleteBy",
ADD COLUMN     "DeletedBy" TEXT;

-- AlterTable
ALTER TABLE "inspection_location" DROP COLUMN "deleteBy",
ADD COLUMN     "DeletedBy" TEXT;

-- AlterTable
ALTER TABLE "license_plate_configuration" DROP COLUMN "deleteBy",
ADD COLUMN     "DeletedBy" TEXT;

-- AlterTable
ALTER TABLE "mobile_device_configuration" DROP COLUMN "deleteBy",
ADD COLUMN     "DeletedBy" TEXT;

-- AlterTable
ALTER TABLE "project" DROP COLUMN "deleteBy",
ADD COLUMN     "DeletedBy" TEXT;

-- AlterTable
ALTER TABLE "project_account" DROP COLUMN "deleteBy",
ADD COLUMN     "DeletedBy" TEXT;

-- AlterTable
ALTER TABLE "vehicle" DROP COLUMN "deleteBy",
ADD COLUMN     "DeletedBy" TEXT;

-- AlterTable
ALTER TABLE "vehicle_image" DROP COLUMN "deleteBy",
ADD COLUMN     "DeletedBy" TEXT;

-- AlterTable
ALTER TABLE "vehicle_imsi" DROP COLUMN "deleteBy",
ADD COLUMN     "DeletedBy" TEXT;

-- AlterTable
ALTER TABLE "vehicle_license_plate" DROP COLUMN "deleteBy",
ADD COLUMN     "DeletedBy" TEXT;
