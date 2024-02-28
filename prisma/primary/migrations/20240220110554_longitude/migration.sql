/*
  Warnings:

  - You are about to drop the column `longtitude` on the `station_location` table. All the data in the column will be lost.
  - Added the required column `longitude` to the `station_location` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "station_location" DROP COLUMN "longtitude",
ADD COLUMN     "longitude" DECIMAL(9,6) NOT NULL;
