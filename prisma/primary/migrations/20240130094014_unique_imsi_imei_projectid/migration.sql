/*
  Warnings:

  - A unique constraint covering the columns `[imei,projectId]` on the table `imei_configuration` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[imsi,projectId]` on the table `imsi_configuration` will be added. If there are existing duplicate values, this will fail.

*/
-- DropIndex
DROP INDEX "imei_configuration_imei_idx";

-- DropIndex
DROP INDEX "imsi_configuration_imsi_key";

-- CreateIndex
CREATE UNIQUE INDEX "imei_configuration_imei_projectId_key" ON "imei_configuration"("imei", "projectId");

-- CreateIndex
CREATE UNIQUE INDEX "imsi_configuration_imsi_projectId_key" ON "imsi_configuration"("imsi", "projectId");
