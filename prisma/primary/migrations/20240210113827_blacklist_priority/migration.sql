/*
  Warnings:

  - You are about to drop the column `priority` on the `imei_configuration` table. All the data in the column will be lost.
  - You are about to drop the column `priority` on the `imsi_configuration` table. All the data in the column will be lost.
  - You are about to drop the column `priority` on the `mobile_device_configuration` table. All the data in the column will be lost.
  - You are about to drop the column `priority` on the `vehicle_license_plate` table. All the data in the column will be lost.

*/
-- AlterTable
ALTER TABLE "imei_configuration" DROP COLUMN "priority",
ADD COLUMN     "blacklistPriority" "BlacklistPriority" NOT NULL DEFAULT 'NORMAL';

-- AlterTable
ALTER TABLE "imsi_configuration" DROP COLUMN "priority",
ADD COLUMN     "blacklistPriority" "BlacklistPriority" NOT NULL DEFAULT 'NORMAL';

-- AlterTable
ALTER TABLE "mobile_device_configuration" DROP COLUMN "priority",
ADD COLUMN     "blacklistPriority" "BlacklistPriority" NOT NULL DEFAULT 'NORMAL';

-- AlterTable
ALTER TABLE "vehicle_license_plate" DROP COLUMN "priority",
ADD COLUMN     "blacklistPriority" "BlacklistPriority" DEFAULT 'NORMAL';
