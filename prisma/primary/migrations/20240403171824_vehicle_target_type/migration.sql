/*
  Warnings:

  - You are about to drop the column `model` on the `vehicle_target_configuration` table. All the data in the column will be lost.
  - You are about to drop the column `modelType` on the `vehicle_target_configuration` table. All the data in the column will be lost.

*/
-- DropForeignKey
ALTER TABLE "vehicle_target_configuration" DROP CONSTRAINT "vehicle_target_configuration_model_modelType_projectId_fkey";

-- AlterTable
ALTER TABLE "vehicle_target_configuration" DROP COLUMN "model",
DROP COLUMN "modelType",
ADD COLUMN     "typeType" "PropertyType" NOT NULL DEFAULT 'VEHICLE_TYPE';

-- AddForeignKey
ALTER TABLE "vehicle_target_configuration" ADD CONSTRAINT "vehicle_target_configuration_type_typeType_projectId_fkey" FOREIGN KEY ("type", "typeType", "projectId") REFERENCES "vehicle_property"("property", "type", "projectId") ON DELETE RESTRICT ON UPDATE CASCADE;
