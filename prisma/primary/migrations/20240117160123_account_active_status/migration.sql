/*
  Warnings:

  - Added the required column `ownerId` to the `project` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "project" ADD COLUMN     "ownerId" UUID NOT NULL;

-- CreateTable
CREATE TABLE "account_configuration" (
    "accountId" UUID NOT NULL,
    "isActive" BOOLEAN NOT NULL DEFAULT true,

    CONSTRAINT "account_configuration_pkey" PRIMARY KEY ("accountId")
);

-- AddForeignKey
ALTER TABLE "account_configuration" ADD CONSTRAINT "account_configuration_accountId_fkey" FOREIGN KEY ("accountId") REFERENCES "account"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "project" ADD CONSTRAINT "project_ownerId_fkey" FOREIGN KEY ("ownerId") REFERENCES "account"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
