/*
  Warnings:

  - You are about to drop the column `ownerId` on the `project` table. All the data in the column will be lost.

*/
-- DropForeignKey
ALTER TABLE "project" DROP CONSTRAINT "project_ownerId_fkey";

-- AlterTable
ALTER TABLE "project" DROP COLUMN "ownerId";
