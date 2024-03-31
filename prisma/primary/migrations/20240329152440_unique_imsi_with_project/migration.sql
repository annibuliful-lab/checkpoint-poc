/*
  Warnings:

  - A unique constraint covering the columns `[imsi,projectId]` on the table `imsi_configuration` will be added. If there are existing duplicate values, this will fail.

*/
-- DropIndex
DROP INDEX "imsi_configuration_imsi_projectId_idx";

-- CreateIndex
CREATE UNIQUE INDEX "imsi_configuration_imsi_projectId_key" ON "imsi_configuration"("imsi", "projectId");
