/*
  Warnings:

  - Added the required column `brand` to the `vehicle_target_configuration` table without a default value. This is not possible if the table is not empty.
  - Added the required column `color` to the `vehicle_target_configuration` table without a default value. This is not possible if the table is not empty.
  - Added the required column `model` to the `vehicle_target_configuration` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "vehicle_target_configuration" ADD COLUMN     "brand" TEXT NOT NULL,
ADD COLUMN     "brandType" "PropertyType" NOT NULL DEFAULT 'VEHICLE_BRAND',
ADD COLUMN     "color" TEXT NOT NULL,
ADD COLUMN     "colorType" "PropertyType" NOT NULL DEFAULT 'VEHICLE_COLOR',
ADD COLUMN     "model" TEXT NOT NULL,
ADD COLUMN     "modelType" "PropertyType" NOT NULL DEFAULT 'VEHICLE_MODEL';

-- AddForeignKey
ALTER TABLE "vehicle_target_configuration" ADD CONSTRAINT "vehicle_target_configuration_brand_brandType_projectId_fkey" FOREIGN KEY ("brand", "brandType", "projectId") REFERENCES "vehicle_property"("property", "type", "projectId") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle_target_configuration" ADD CONSTRAINT "vehicle_target_configuration_color_colorType_projectId_fkey" FOREIGN KEY ("color", "colorType", "projectId") REFERENCES "vehicle_property"("property", "type", "projectId") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vehicle_target_configuration" ADD CONSTRAINT "vehicle_target_configuration_model_modelType_projectId_fkey" FOREIGN KEY ("model", "modelType", "projectId") REFERENCES "vehicle_property"("property", "type", "projectId") ON DELETE RESTRICT ON UPDATE CASCADE;
