import { DevicePermittedLabel } from "@/apollo-client";
import { VehicleTransection } from "./types";

export const TABLE_HEAD = [
  { id: "title", label: "Arrival Time", minWidth: 160 },
  { id: "title", label: "License Plate", minWidth: 140 },
  { id: "title", label: "License Plate Type", minWidth: 200 },
  { id: "title", label: "Brand", minWidth: 116 },
  { id: "title", label: "Vehicle Type", minWidth: 140 },
  { id: "title", label: "Color", minWidth: 116 },
  { id: "title", label: "Station Site", minWidth: 140 },
  { id: "title", label: "IMEI amount", minWidth: 140, align: "center" },
  { id: "title", label: "IMSI amount ", minWidth: 140, align: "center" },
  { id: "title", label: "Tags", minWidth: 160 },
  { id: "title", label: "Remark", minWidth: 160 },
];
export const VEHICLR_TRANSECTIONS: VehicleTransection[] = [];
