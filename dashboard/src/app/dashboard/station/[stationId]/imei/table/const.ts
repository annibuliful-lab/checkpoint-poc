import { DevicePermittedLabel } from "@/apollo-client";
import { ImeiImsiTransaction } from "./types";

export const TABLE_HEAD = [
  { id: "title", label: "Arrival Time", width: 116 },
  { id: "title", label: "IMEI Number", width: 116 },
  { id: "title", label: "IMSI Number", width: 116 },
  { id: "title", label: "Phone model", width: 116 },
  { id: "title", label: "License Plate", width: 116 },
  { id: "title", label: "Station Site", width: 116 },
  { id: "title", label: "Tags", width: 116 },
];
export const IMEI_IMSI_TRANSECTIONS: ImeiImsiTransaction[] = [];
