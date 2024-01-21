/*
  Warnings:

  - A unique constraint covering the columns `[subject,action]` on the table `permission` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[title]` on the table `project` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[projectId,title]` on the table `project_role` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[roleId,permissionId]` on the table `project_role_permission` will be added. If there are existing duplicate values, this will fail.
  - Added the required column `projectId` to the `project_account` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "project_account" ADD COLUMN     "projectId" UUID NOT NULL;

-- CreateIndex
CREATE UNIQUE INDEX "permission_subject_action_key" ON "permission"("subject", "action");

-- CreateIndex
CREATE UNIQUE INDEX "project_title_key" ON "project"("title");

-- CreateIndex
CREATE UNIQUE INDEX "project_role_projectId_title_key" ON "project_role"("projectId", "title");

-- CreateIndex
CREATE UNIQUE INDEX "project_role_permission_roleId_permissionId_key" ON "project_role_permission"("roleId", "permissionId");

-- AddForeignKey
ALTER TABLE "project_account" ADD CONSTRAINT "project_account_projectId_fkey" FOREIGN KEY ("projectId") REFERENCES "project"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
