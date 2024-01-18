-- AlterTable
ALTER TABLE "session_token" ADD COLUMN     "isRefreshToken" BOOLEAN NOT NULL DEFAULT false;
