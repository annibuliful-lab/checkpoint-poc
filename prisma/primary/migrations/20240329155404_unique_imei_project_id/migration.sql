/*
  Warnings:

  - A unique constraint covering the columns `[imei,projectId]` on the table `imei_configuration` will be added. If there are existing duplicate values, this will fail.

*/
-- DropIndex
DROP INDEX "imei_configuration_imei_projectId_idx";

-- CreateIndex
CREATE UNIQUE INDEX "imei_configuration_imei_projectId_key" ON "imei_configuration"("imei", "projectId");
