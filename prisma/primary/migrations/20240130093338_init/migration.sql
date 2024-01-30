CREATE EXTENSION pg_trgm;

-- CreateEnum
CREATE TYPE "PermissionAction" AS ENUM ('CREATE', 'UPDATE', 'DELETE', 'READ');

-- CreateEnum
CREATE TYPE "DeviceCaptureStatus" AS ENUM ('WHITELIST', 'BLACKLIST');

-- CreateEnum
CREATE TYPE "RemarkState" AS ENUM ('WHITELIST', 'BLACKLIST', 'IN_QUEUE', 'IN_PROGRESS', 'PASSED', 'WAITING', 'INVESTIGATING', 'SUSPICION');

-- CreateTable
CREATE TABLE "session_token" (
    "token" TEXT NOT NULL,
    "revoke" BOOLEAN NOT NULL DEFAULT false,
    "isRefreshToken" BOOLEAN NOT NULL DEFAULT false,
    "accountId" UUID NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3),

    CONSTRAINT "session_token_pkey" PRIMARY KEY ("token")
);

-- CreateTable
CREATE TABLE "permission" (
    "id" UUID NOT NULL,
    "subject" TEXT NOT NULL,
    "action" "PermissionAction" NOT NULL,

    CONSTRAINT "permission_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "account" (
    "id" UUID NOT NULL,
    "username" VARCHAR(24) NOT NULL,
    "password" TEXT NOT NULL,
    "createdBy" TEXT NOT NULL,
    "updatedBy" TEXT,
    "deleteBy" TEXT,
    "deletedAt" TIMESTAMP(3),
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3),

    CONSTRAINT "account_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "account_configuration" (
    "accountId" UUID NOT NULL,
    "isActive" BOOLEAN NOT NULL DEFAULT true,

    CONSTRAINT "account_configuration_pkey" PRIMARY KEY ("accountId")
);

-- CreateTable
CREATE TABLE "project" (
    "id" UUID NOT NULL,
    "title" TEXT NOT NULL,
    "createdBy" TEXT NOT NULL,
    "updatedBy" TEXT,
    "deleteBy" TEXT,
    "deletedAt" TIMESTAMP(3),
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3),

    CONSTRAINT "project_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "project_role" (
    "id" UUID NOT NULL,
    "projectId" UUID NOT NULL,
    "title" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP(3),

    CONSTRAINT "project_role_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "project_role_permission" (
    "id" UUID NOT NULL,
    "projectId" UUID NOT NULL,
    "roleId" UUID NOT NULL,
    "permissionId" UUID NOT NULL,

    CONSTRAINT "project_role_permission_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "project_account" (
    "id" UUID NOT NULL,
    "projectId" UUID NOT NULL,
    "accountId" UUID NOT NULL,
    "roleId" UUID NOT NULL,
    "isActive" BOOLEAN NOT NULL DEFAULT true,
    "createdBy" TEXT NOT NULL,
    "updatedBy" TEXT,
    "deleteBy" TEXT,
    "deletedAt" TIMESTAMP(3),
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3),

    CONSTRAINT "project_account_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "imsi_configuration" (
    "id" UUID NOT NULL,
    "projectId" UUID NOT NULL,
    "imsi" VARCHAR(24) NOT NULL,
    "status" "DeviceCaptureStatus" NOT NULL,
    "createdBy" TEXT NOT NULL,
    "updatedBy" TEXT,
    "deleteBy" TEXT,
    "deletedAt" TIMESTAMP(3),
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3),

    CONSTRAINT "imsi_configuration_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "imei_configuration" (
    "id" UUID NOT NULL,
    "projectId" UUID NOT NULL,
    "imei" VARCHAR(24) NOT NULL,
    "status" "DeviceCaptureStatus" NOT NULL,
    "createdBy" TEXT NOT NULL,
    "updatedBy" TEXT,
    "deleteBy" TEXT,
    "deletedAt" TIMESTAMP(3),
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3),

    CONSTRAINT "imei_configuration_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "mobile_device_configuration" (
    "id" UUID NOT NULL,
    "projectId" UUID NOT NULL,
    "title" TEXT NOT NULL,
    "tags" TEXT[],
    "imsi" TEXT NOT NULL,
    "imei" TEXT NOT NULL,
    "msisdn" TEXT,
    "status" "DeviceCaptureStatus" NOT NULL,
    "createdBy" TEXT NOT NULL,
    "updatedBy" TEXT,
    "deleteBy" TEXT,
    "deletedAt" TIMESTAMP(3),
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3),

    CONSTRAINT "mobile_device_configuration_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "inspection_location" (
    "id" UUID NOT NULL,
    "projectId" UUID NOT NULL,
    "title" TEXT NOT NULL,
    "description" TEXT,
    "tags" TEXT[] DEFAULT ARRAY[]::TEXT[],
    "latitude" DECIMAL(8,6) NOT NULL,
    "longtitude" DECIMAL(9,6) NOT NULL,
    "createdBy" TEXT NOT NULL,
    "updatedBy" TEXT,
    "deleteBy" TEXT,
    "deletedAt" TIMESTAMP(3),
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3),

    CONSTRAINT "inspection_location_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "license_plate_configuration" (
    "id" UUID NOT NULL,
    "projectId" UUID NOT NULL,
    "prefix" VARCHAR(4) NOT NULL,
    "number" VARCHAR(24) NOT NULL,
    "province" TEXT NOT NULL,
    "country" TEXT,
    "status" "DeviceCaptureStatus" NOT NULL,
    "createdBy" TEXT NOT NULL,
    "updatedBy" TEXT,
    "deleteBy" TEXT,
    "deletedAt" TIMESTAMP(3),
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3),

    CONSTRAINT "license_plate_configuration_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "vehicle" (
    "id" UUID NOT NULL,
    "projectId" UUID NOT NULL,
    "inspectionLocationId" UUID NOT NULL,
    "brand" TEXT NOT NULL,
    "color" TEXT NOT NULL,
    "model" TEXT NOT NULL,
    "status" "RemarkState" NOT NULL DEFAULT 'IN_PROGRESS',
    "createdBy" TEXT NOT NULL,
    "updatedBy" TEXT,
    "deleteBy" TEXT,
    "deletedAt" TIMESTAMP(3),
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3),

    CONSTRAINT "vehicle_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "vehicle_license_plate" (
    "vehicleId" UUID NOT NULL,
    "prefix" VARCHAR(4) NOT NULL,
    "number" VARCHAR(8) NOT NULL,
    "province" VARCHAR(32) NOT NULL,
    "country" TEXT,
    "status" "DeviceCaptureStatus",
    "createdBy" TEXT NOT NULL,
    "updatedBy" TEXT,
    "deleteBy" TEXT,
    "deletedAt" TIMESTAMP(3),
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3),

    CONSTRAINT "vehicle_license_plate_pkey" PRIMARY KEY ("vehicleId")
);

-- CreateTable
CREATE TABLE "vehicle_image" (
    "id" UUID NOT NULL,
    "vehicleId" UUID NOT NULL,
    "s3Key" TEXT NOT NULL,
    "status" "DeviceCaptureStatus",
    "createdBy" TEXT NOT NULL,
    "updatedBy" TEXT,
    "deleteBy" TEXT,
    "deletedAt" TIMESTAMP(3),
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3),

    CONSTRAINT "vehicle_image_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "vehicle_imsi" (
    "id" UUID NOT NULL,
    "vehicleId" UUID NOT NULL,
    "labelStatus" "RemarkState" NOT NULL DEFAULT 'IN_PROGRESS',
    "imsi" VARCHAR(24) NOT NULL,
    "imei" VARCHAR(24) NOT NULL,
    "status" "DeviceCaptureStatus",
    "signalStrength" DOUBLE PRECISION NOT NULL,
    "registerDateTime" TIMESTAMP(3) NOT NULL,
    "createdBy" TEXT NOT NULL,
    "updatedBy" TEXT,
    "deleteBy" TEXT,
    "deletedAt" TIMESTAMP(3),
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3),

    CONSTRAINT "vehicle_imsi_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "permission_subject_action_key" ON "permission"("subject", "action");

-- CreateIndex
CREATE UNIQUE INDEX "account_username_key" ON "account"("username");

-- CreateIndex
CREATE UNIQUE INDEX "project_title_key" ON "project"("title");

-- CreateIndex
CREATE UNIQUE INDEX "project_role_projectId_title_key" ON "project_role"("projectId", "title");

-- CreateIndex
CREATE UNIQUE INDEX "project_role_permission_roleId_permissionId_key" ON "project_role_permission"("roleId", "permissionId");

-- CreateIndex
CREATE UNIQUE INDEX "imsi_configuration_imsi_key" ON "imsi_configuration"("imsi");

-- CreateIndex
CREATE INDEX "imei_configuration_imei_idx" ON "imei_configuration"("imei");

-- CreateIndex
CREATE INDEX "mobile_device_configuration_tags_idx" ON "mobile_device_configuration" USING GIN ("tags");

-- CreateIndex
CREATE INDEX "mobile_device_configuration_title_idx" ON "mobile_device_configuration" USING GIN ("title" gin_trgm_ops);

-- CreateIndex
CREATE INDEX "inspection_location_tags_idx" ON "inspection_location" USING GIN ("tags");

-- CreateIndex
CREATE INDEX "inspection_location_title_idx" ON "inspection_location" USING GIN ("title" gin_trgm_ops);

-- CreateIndex
CREATE INDEX "license_plate_configuration_prefix_idx" ON "license_plate_configuration"("prefix");

-- CreateIndex
CREATE INDEX "license_plate_configuration_number_prefix_idx" ON "license_plate_configuration"("number", "prefix");

-- CreateIndex
CREATE INDEX "license_plate_configuration_number_prefix_province_idx" ON "license_plate_configuration"("number", "prefix", "province");

-- CreateIndex
CREATE INDEX "vehicle_brand_idx" ON "vehicle"("brand");

-- CreateIndex
CREATE INDEX "vehicle_color_idx" ON "vehicle"("color");

-- CreateIndex
CREATE INDEX "vehicle_model_idx" ON "vehicle"("model");

-- CreateIndex
CREATE INDEX "vehicle_license_plate_prefix_idx" ON "vehicle_license_plate"("prefix");

-- CreateIndex
CREATE INDEX "vehicle_license_plate_number_prefix_idx" ON "vehicle_license_plate"("number", "prefix");

-- CreateIndex
CREATE INDEX "vehicle_license_plate_number_prefix_province_idx" ON "vehicle_license_plate"("number", "prefix", "province");

-- CreateIndex
CREATE INDEX "vehicle_imsi_imei_imsi_idx" ON "vehicle_imsi"("imei", "imsi");

-- AddForeignKey
ALTER TABLE "session_token" ADD CONSTRAINT "session_token_accountId_fkey" FOREIGN KEY ("accountId") REFERENCES "account"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "account_configuration" ADD CONSTRAINT "account_configuration_accountId_fkey" FOREIGN KEY ("accountId") REFERENCES "account"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "project_role" ADD CONSTRAINT "project_role_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES "project"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "project_role_permission" ADD CONSTRAINT "project_role_permission_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES "project"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "project_role_permission" ADD CONSTRAINT "project_role_permission_roleId_fkey" FOREIGN KEY ("roleId") REFERENCES "project_role"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "project_role_permission" ADD CONSTRAINT "project_role_permission_permissionId_fkey" FOREIGN KEY ("permissionId") REFERENCES "permission"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "project_account" ADD CONSTRAINT "project_account_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES "project"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "project_account" ADD CONSTRAINT "project_account_accountId_fkey" FOREIGN KEY ("accountId") REFERENCES "account"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "project_account" ADD CONSTRAINT "project_account_roleId_fkey" FOREIGN KEY ("roleId") REFERENCES "project_role"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "imsi_configuration" ADD CONSTRAINT "imsi_configuration_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES "project"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "imei_configuration" ADD CONSTRAINT "imei_configuration_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES "project"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "mobile_device_configuration" ADD CONSTRAINT "mobile_device_configuration_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES "project"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "inspection_location" ADD CONSTRAINT "inspection_location_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES "project"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "license_plate_configuration" ADD CONSTRAINT "license_plate_configuration_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES "project"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle" ADD CONSTRAINT "vehicle_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES "project"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle" ADD CONSTRAINT "vehicle_inspectionLocationId_fkey" FOREIGN KEY ("inspectionLocationId") REFERENCES "inspection_location"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle_license_plate" ADD CONSTRAINT "vehicle_license_plate_vehicleId_fkey" FOREIGN KEY ("vehicleId") REFERENCES "vehicle"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle_image" ADD CONSTRAINT "vehicle_image_vehicleId_fkey" FOREIGN KEY ("vehicleId") REFERENCES "vehicle"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle_imsi" ADD CONSTRAINT "vehicle_imsi_vehicleId_fkey" FOREIGN KEY ("vehicleId") REFERENCES "vehicle"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
