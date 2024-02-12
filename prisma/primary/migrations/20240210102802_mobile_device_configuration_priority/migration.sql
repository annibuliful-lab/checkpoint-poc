/*
  Warnings:

  - You are about to drop the column `label` on the `mobile_device_configuration` table. All the data in the column will be lost.

*/
-- AlterTable
ALTER TABLE "mobile_device_configuration" DROP COLUMN "label",
ADD COLUMN     "permittedLabel" "DevicePermittedLabel" NOT NULL DEFAULT 'NONE',
ADD COLUMN     "priority" "BlacklistPriority" NOT NULL DEFAULT 'NORMAL';
