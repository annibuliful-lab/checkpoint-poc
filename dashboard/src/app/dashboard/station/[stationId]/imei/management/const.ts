import { BlacklistPriority, DevicePermittedLabel } from "@/apollo-client";

export const TABLE_HEAD = [
  { id: "title", label: "IMEI Number", width: 116 },
  { id: "status", label: "IMSI Number", width: 116 },
  { id: "updatedAt", label: "License Plate", width: 116 },
  { id: "geolocation", label: "Status", width: 116 },
  { id: "tags", label: "Tag", width: 116 },
  { id: "action", label: "Action", width: 116 },
];
export const IMEI_STATUS_OPTIONS = Object.keys(DevicePermittedLabel).map(
  (key) => ({ id: DevicePermittedLabel[key as "Blacklist"], label: key })
);
export const IMEI_PRIORITY_OPTIONS = Object.keys(BlacklistPriority).map(
  (key) => ({ id: BlacklistPriority[key as "Normal"], label: key })
);
