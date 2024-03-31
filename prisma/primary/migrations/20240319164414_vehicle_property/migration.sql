-- CreateEnum
CREATE TYPE "PropertyType" AS ENUM ('VEHICLE_COLOR', 'LP_TYPE', 'VEHICLE_BRAND', 'VEHICLE_TYPE', 'VEHICLE_MODEL');

-- DropIndex
DROP INDEX "station_vehicle_activity_brand_idx";

-- DropIndex
DROP INDEX "station_vehicle_activity_color_idx";

-- DropIndex
DROP INDEX "station_vehicle_activity_model_idx";

-- AlterTable
ALTER TABLE "station_vehicle_activity" ADD COLUMN     "brandType" "PropertyType" NOT NULL DEFAULT 'VEHICLE_BRAND',
ADD COLUMN     "colorType" "PropertyType" NOT NULL DEFAULT 'VEHICLE_COLOR',
ADD COLUMN     "modelType" "PropertyType" NOT NULL DEFAULT 'VEHICLE_MODEL';

-- CreateTable
CREATE TABLE "vehicle_property" (
    "projectId" UUID NOT NULL,
    "property" TEXT NOT NULL,
    "type" "PropertyType" NOT NULL,

    CONSTRAINT "vehicle_property_pkey" PRIMARY KEY ("property","type","projectId")
);

-- AddForeignKey
ALTER TABLE "vehicle_property" ADD CONSTRAINT "vehicle_property_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES "project"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "station_vehicle_activity" ADD CONSTRAINT "station_vehicle_activity_brand_brandType_projectId_fkey" FOREIGN KEY ("brand", "brandType", "projectId") REFERENCES "vehicle_property"("property", "type", "projectId") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "station_vehicle_activity" ADD CONSTRAINT "station_vehicle_activity_color_colorType_projectId_fkey" FOREIGN KEY ("color", "colorType", "projectId") REFERENCES "vehicle_property"("property", "type", "projectId") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "station_vehicle_activity" ADD CONSTRAINT "station_vehicle_activity_model_modelType_projectId_fkey" FOREIGN KEY ("model", "modelType", "projectId") REFERENCES "vehicle_property"("property", "type", "projectId") ON DELETE RESTRICT ON UPDATE CASCADE;
