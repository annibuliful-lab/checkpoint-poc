-- DropIndex
DROP INDEX "imei_configuration_imei_projectId_key";

-- DropIndex
DROP INDEX "imsi_configuration_imsi_projectId_key";

-- CreateIndex
CREATE INDEX "imei_configuration_imei_projectId_idx" ON "imei_configuration"("imei", "projectId");

-- CreateIndex
CREATE INDEX "imsi_configuration_imsi_projectId_idx" ON "imsi_configuration"("imsi", "projectId");
