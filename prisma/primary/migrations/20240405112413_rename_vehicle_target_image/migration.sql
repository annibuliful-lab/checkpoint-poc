/*
  Warnings:

  - You are about to drop the `vehicle_target_configurationimage` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropForeignKey
ALTER TABLE "vehicle_target_configurationimage" DROP CONSTRAINT "vehicle_target_configurationimage_vehicleTargetConfigurati_fkey";

-- DropTable
DROP TABLE "vehicle_target_configurationimage";

-- CreateTable
CREATE TABLE "vehicle_target_configuration_image" (
    "id" UUID NOT NULL,
    "vehicleTargetConfigurationId" UUID NOT NULL,
    "type" "ImageType" NOT NULL,
    "s3Key" TEXT NOT NULL,
    "createdBy" TEXT NOT NULL,
    "updatedBy" TEXT,
    "deletedBy" TEXT,
    "deletedAt" TIMESTAMPTZ(3),
    "createdAt" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMPTZ(3),

    CONSTRAINT "vehicle_target_configuration_image_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "vehicle_target_configuration_image" ADD CONSTRAINT "vehicle_target_configuration_image_vehicleTargetConfigurat_fkey" FOREIGN KEY ("vehicleTargetConfigurationId") REFERENCES "vehicle_target_configuration"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
