/*
  Warnings:

  - You are about to drop the column `status` on the `imei_configuration` table. All the data in the column will be lost.
  - You are about to drop the column `status` on the `imsi_configuration` table. All the data in the column will be lost.
  - You are about to drop the column `status` on the `license_plate_configuration` table. All the data in the column will be lost.
  - You are about to drop the column `status` on the `mobile_device_configuration` table. All the data in the column will be lost.
  - You are about to drop the column `status` on the `vehicle_image` table. All the data in the column will be lost.
  - You are about to drop the column `labelStatus` on the `vehicle_imsi` table. All the data in the column will be lost.
  - You are about to drop the column `status` on the `vehicle_imsi` table. All the data in the column will be lost.
  - You are about to drop the column `status` on the `vehicle_license_plate` table. All the data in the column will be lost.
  - Added the required column `label` to the `imei_configuration` table without a default value. This is not possible if the table is not empty.
  - Added the required column `label` to the `imsi_configuration` table without a default value. This is not possible if the table is not empty.
  - Added the required column `label` to the `license_plate_configuration` table without a default value. This is not possible if the table is not empty.
  - Added the required column `label` to the `mobile_device_configuration` table without a default value. This is not possible if the table is not empty.
  - Added the required column `referenceImeiConfigurationId` to the `mobile_device_configuration` table without a default value. This is not possible if the table is not empty.
  - Added the required column `referenceImsiConfigurationId` to the `mobile_device_configuration` table without a default value. This is not possible if the table is not empty.

*/
-- CreateEnum
CREATE TYPE "DeviceCaptureLabel" AS ENUM ('WHITELIST', 'BLACKLIST');

-- AlterTable
ALTER TABLE "imei_configuration" DROP COLUMN "status",
ADD COLUMN     "label" "DeviceCaptureLabel" NOT NULL;

-- AlterTable
ALTER TABLE "imsi_configuration" DROP COLUMN "status",
ADD COLUMN     "label" "DeviceCaptureLabel" NOT NULL;

-- AlterTable
ALTER TABLE "license_plate_configuration" DROP COLUMN "status",
ADD COLUMN     "label" "DeviceCaptureLabel" NOT NULL;

-- AlterTable
ALTER TABLE "mobile_device_configuration" DROP COLUMN "status",
ADD COLUMN     "label" "DeviceCaptureLabel" NOT NULL,
ADD COLUMN     "referenceImeiConfigurationId" UUID NOT NULL,
ADD COLUMN     "referenceImsiConfigurationId" UUID NOT NULL;

-- AlterTable
ALTER TABLE "vehicle" ALTER COLUMN "status" DROP NOT NULL,
ALTER COLUMN "status" DROP DEFAULT;

-- AlterTable
ALTER TABLE "vehicle_image" DROP COLUMN "status",
ADD COLUMN     "label" "DeviceCaptureLabel";

-- AlterTable
ALTER TABLE "vehicle_imsi" DROP COLUMN "labelStatus",
DROP COLUMN "status",
ADD COLUMN     "label" "DeviceCaptureLabel",
ADD COLUMN     "remarkStatus" "RemarkState" NOT NULL DEFAULT 'IN_PROGRESS';

-- AlterTable
ALTER TABLE "vehicle_license_plate" DROP COLUMN "status",
ADD COLUMN     "label" "DeviceCaptureLabel";

-- DropEnum
DROP TYPE "DeviceCaptureStatus";

-- AddForeignKey
ALTER TABLE "mobile_device_configuration" ADD CONSTRAINT "mobile_device_configuration_referenceImsiConfigurationId_fkey" FOREIGN KEY ("referenceImsiConfigurationId") REFERENCES "imsi_configuration"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "mobile_device_configuration" ADD CONSTRAINT "mobile_device_configuration_referenceImeiConfigurationId_fkey" FOREIGN KEY ("referenceImeiConfigurationId") REFERENCES "imei_configuration"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
