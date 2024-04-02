import { DevicePermittedLabel } from "@/apollo-client";
import { ImsiImeiTransaction } from "./types";

export const TABLE_HEAD = [
  { id: "title", label: "Arrival Time", width: 116 },
  { id: "title", label: "IMEI Number", width: 116 },
  { id: "title", label: "IMSI Number", width: 116 },
  { id: "title", label: "Phone model", width: 116 },
  { id: "title", label: "Tags", width: 116 },
];
export const DATA: ImsiImeiTransaction[] = [
  {
    id: "1",
    arrivalTime: new Date("2024-03-27T08:00:00Z"),
    imei: { status: DevicePermittedLabel.None, imei: 123456789012345 },
    imsi: { status: DevicePermittedLabel.None, imsi: 987654321098765 },
    phoneModel: "ExamplePhone1",
    tags: ["tag1", "tag2"],
  },
  {
    id: "2",
    arrivalTime: new Date("2024-03-27T09:30:00Z"),
    imei: { status: DevicePermittedLabel.Blacklist, imei: 234567890123456 },
    imsi: { status: DevicePermittedLabel.None, imsi: 876543210987654 },
    phoneModel: "ExamplePhone2",
    tags: ["tag2", "tag3"],
  },
  {
    id: "3",
    arrivalTime: new Date("2024-03-27T11:15:00Z"),
    imei: { status: DevicePermittedLabel.None, imei: 345678901234567 },
    imsi: { status: DevicePermittedLabel.Whitelist, imsi: 765432109876543 },
    phoneModel: "ExamplePhone3",
    tags: ["tag3", "tag4"],
  },
];
