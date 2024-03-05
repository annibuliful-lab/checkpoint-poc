-- AlterTable
ALTER TABLE "imei_configuration" ALTER COLUMN "blacklistPriority" SET DEFAULT 'NONE';

-- AlterTable
ALTER TABLE "imsi_configuration" ALTER COLUMN "blacklistPriority" SET DEFAULT 'NONE';

-- AlterTable
ALTER TABLE "license_plate_configuration" ALTER COLUMN "blacklistPriority" SET DEFAULT 'NONE';

-- AlterTable
ALTER TABLE "mobile_device_configuration" ALTER COLUMN "blacklistPriority" SET DEFAULT 'NONE';

-- AlterTable
ALTER TABLE "vehicle_license_plate" ALTER COLUMN "blacklistPriority" SET DEFAULT 'NONE';
