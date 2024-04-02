import { Authentication } from "@/apollo-client";
import { atomWithStorage } from "jotai/utils";
export const authAtom = atomWithStorage<
  (Authentication & { projectId: string }) | undefined
>("auth", undefined);
