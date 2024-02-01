/*
  Warnings:

  - You are about to drop the column `imei` on the `mobile_device_configuration` table. All the data in the column will be lost.
  - You are about to drop the column `imsi` on the `mobile_device_configuration` table. All the data in the column will be lost.

*/
-- AlterTable
ALTER TABLE "mobile_device_configuration" DROP COLUMN "imei",
DROP COLUMN "imsi";
