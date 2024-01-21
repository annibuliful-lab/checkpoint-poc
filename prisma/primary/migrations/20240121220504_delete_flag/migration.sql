-- AlterTable
ALTER TABLE "project" ADD COLUMN     "deleteBy" TEXT,
ADD COLUMN     "deletedAt" TIMESTAMP(3);

-- AlterTable
ALTER TABLE "project_role" ADD COLUMN     "deletedAt" TIMESTAMP(3);
