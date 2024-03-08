/*
  Warnings:

  - You are about to drop the `license_plate_configuration` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `license_plate_configuration_tag` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropForeignKey
ALTER TABLE "license_plate_configuration" DROP CONSTRAINT "license_plate_configuration_projectId_fkey";

-- DropForeignKey
ALTER TABLE "license_plate_configuration_tag" DROP CONSTRAINT "license_plate_configuration_tag_licensePlateConfigurationI_fkey";

-- DropForeignKey
ALTER TABLE "license_plate_configuration_tag" DROP CONSTRAINT "license_plate_configuration_tag_tagId_fkey";

-- DropTable
DROP TABLE "license_plate_configuration";

-- DropTable
DROP TABLE "license_plate_configuration_tag";

-- CreateTable
CREATE TABLE "vehicle_target_configuration_tag" (
    "id" UUID NOT NULL,
    "createdAt" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "vehicleTargetConfigurationId" UUID NOT NULL,
    "tagId" UUID NOT NULL,
    "createdBy" TEXT NOT NULL,

    CONSTRAINT "vehicle_target_configuration_tag_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "vehicle_target_configuration" (
    "id" UUID NOT NULL,
    "projectId" UUID NOT NULL,
    "prefix" VARCHAR(4) NOT NULL,
    "number" VARCHAR(24) NOT NULL,
    "province" TEXT NOT NULL,
    "type" TEXT NOT NULL,
    "country" TEXT,
    "permittedLabel" "DevicePermittedLabel" NOT NULL DEFAULT 'NONE',
    "blacklistPriority" "BlacklistPriority" NOT NULL DEFAULT 'NONE',
    "createdBy" TEXT NOT NULL,
    "updatedBy" TEXT,
    "deletedBy" TEXT,
    "deletedAt" TIMESTAMPTZ(3),
    "createdAt" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMPTZ(3),

    CONSTRAINT "vehicle_target_configuration_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE INDEX "vehicle_target_configuration_prefix_idx" ON "vehicle_target_configuration"("prefix");

-- CreateIndex
CREATE INDEX "vehicle_target_configuration_number_prefix_idx" ON "vehicle_target_configuration"("number", "prefix");

-- CreateIndex
CREATE INDEX "vehicle_target_configuration_number_prefix_province_idx" ON "vehicle_target_configuration"("number", "prefix", "province");

-- AddForeignKey
ALTER TABLE "vehicle_target_configuration_tag" ADD CONSTRAINT "vehicle_target_configuration_tag_vehicleTargetConfiguratio_fkey" FOREIGN KEY ("vehicleTargetConfigurationId") REFERENCES "vehicle_target_configuration"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle_target_configuration_tag" ADD CONSTRAINT "vehicle_target_configuration_tag_tagId_fkey" FOREIGN KEY ("tagId") REFERENCES "tag"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle_target_configuration" ADD CONSTRAINT "vehicle_target_configuration_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES "project"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
