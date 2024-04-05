-- CreateTable
CREATE TABLE "vehicle_target_configurationimage" (
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

    CONSTRAINT "vehicle_target_configurationimage_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "vehicle_target_configurationimage" ADD CONSTRAINT "vehicle_target_configurationimage_vehicleTargetConfigurati_fkey" FOREIGN KEY ("vehicleTargetConfigurationId") REFERENCES "vehicle_target_configuration"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
