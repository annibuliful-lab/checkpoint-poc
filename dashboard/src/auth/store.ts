import { Authentication } from "@/apollo-client";
import { atomWithStorage } from "jotai/utils";
export const authAtom = atomWithStorage<Authentication | undefined>(
  "auth",
  undefined
);
