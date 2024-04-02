import { StationLocation } from "@/apollo-client";
import { atomWithStorage } from "jotai/utils";

export const stationHistoriesAtom = atomWithStorage<
  Record<string, StationLocation>
>("station-histories", {});
