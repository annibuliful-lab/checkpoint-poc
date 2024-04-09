import { StationDashboardActivity } from "@/apollo-client";
import { atom } from "jotai";
export const stationDashboardActivityAtom = atom<
  StationDashboardActivity | undefined
>(undefined);
